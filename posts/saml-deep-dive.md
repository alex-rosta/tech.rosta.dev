---
title: "Deep Dive into SAML"
tags: "saml, authentication, identity"
created: 2025-07-23
updated: 2025-07-23
---
In this post, we'll explore the Security Assertion Markup Language (SAML), a tool for single sign-on (SSO). SAML is used to faciliate trust between an application and an identity provider (IdP), allowing users to authenticate once and gain access to multiple applications without needing to log in again.

![image](https://cdn.prod.website-files.com/638117e276a1ed457b80b0fe/65c25213371eaef64ed6d0cb_saml-2-logo.png)

## The Trust Relationship

Trust is established between the service provider (SP) and the identity provider (IdP) through metadata exchange.

The SP provides its metadata to the IdP, which includes information like the SP's entity ID and the endpoints for SAML requests and responses. The IdP also provides its metadata to the SP, which includes the IdP's entity ID and public key for signature verification.

Example IdP metadata:

```xml
<EntityDescriptor
    ID="_c066524f-ba36-49d5-9dfa-ae14e13c1392"
    entityID="https://idp.com"
    xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
    xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">
    <IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.com/saml/sso" />
        <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://idp.com/saml/sso" />
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.com/saml/slo" />
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://idp.com/saml/slo" />
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Artifact" Location="https://idp.com/saml/slo" />
        <KeyDescriptor use="signing">
            <KeyInfo
                xmlns="http://www.w3.org/2000/09/xmldsig#">
                <X509Data>
                    <X509Certificate>IDP_PUBLIC_SIGNING_CERTIFICATE_USED_FOR_SIGNING_RESPONSES</X509Certificate>
                </X509Data>
            </KeyInfo>
        </KeyDescriptor>
    </IDPSSODescriptor>
</EntityDescriptor>
```

Important elements in the IdP metadata include:

- `SingleSignOnService`: The endpoint where the SP sends authentication requests. It can support multiple bindings (HTTP-Redirect, HTTP-POST).
- `entityID`: The unique identifier for the IdP. SP can use this to identify the IdP, as they may have multiple IdPs.
- `KeyDescriptor`: Contains the public key used by the IdP to sign SAML responses. The SP uses this key to verify the authenticity of the response. Often this is self-signed, so the SP must trust the IdP's certificate.
- `SingleLogoutService`: The endpoint for logging out of the IdP.

In my opinion the IdP should never advertise claims / nameid-format or attributes in the metadata, as that can (and should) change depending on which SP is requesting a ticket.
Unless the IdP is configured to always return the same claims for all SPs, which is unusual.

Example SP metadata:

```xml
<EntityDescriptor 
    ID="_33fc2606-cc4f-4883-b3a4-2c1d37090848"
    entityID="https://sp.com"
    xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
    xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">
    <SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://sp.com/saml/acs" index="1" />
        <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://sp.com/saml/acs" index="2" />
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://sp.com/saml/slo" />
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://sp.com/saml/slo" />
        <KeyDescriptor use="signing">
            <KeyInfo
                xmlns="http://www.w3.org/2000/09/xmldsig#">
                <X509Data>
                    <X509Certificate>SP_PUBLIC_SIGNING_CERTIFICATE_USED_FOR_SIGNING_REQUESTS</X509Certificate>
                </X509Data>
            </KeyInfo>
        </KeyDescriptor>
    </SPSSODescriptor>
</EntityDescriptor>
```

Important elements in the SP metadata include:

- `entityID`: The unique identifier for the SP. The IdP uses this to identify the SP, as they are likely to have multiple SPs.
- `AssertionConsumerService`: The endpoint where the IdP sends SAML assertions after successful authentication
- `KeyDescriptor`: Contains the public key used by the SP to sign SAML requests. The IdP uses this key to verify the authenticity of the request.
- `SingleLogoutService`: The endpoint for logging out of the SP.

As you see they are quite similar, but the SP metadata does not include the `SingleSignOnService` element, as that is only relevant for the IdP.
Sometimes the SP may advertise an encryption certificate, but I feel this is overkill in most cases. And also extremely hard to troubleshoot when things go wrong with assertions as they are encrypted...

## Let's Walk Through a SAML Flow

User accesses an application (SP) that requires authentication. The SP checks if the user is authenticated. If not, it initiates a SAML authentication request.

Example SAML authentication request:

```xml
<samlp:AuthnRequest
  xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
  ID="C2dE3fH4iJ5kL6mN7oP8qR9sT0uV1w"
  Version="2.0" IssueInstant="2013-03-18T03:28:54.1839884Z"
  xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol">
  <Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion">https://sp.com</Issuer>
</samlp:AuthnRequest>
```

IdP receives the request and takes note of the `Issuer` element, which identifies the SP. User logs in using one of the authentication methods. Let's say [Freja eID](https://frejaeid.com/) in this case.

Freja eID returns multiple attributes about the user, such as full name user personal number etc.
The IdP creates a SAML assertion containing the user's identity information and signs it with its private key.

Before the assertion, the IdP may also perform lookups to third party catalogues to enrich the user's attributes, such as fetching the user's internal account name or group memberships in Active Directory, information it did not have after the authentication with Freja. But this is not a requirement of SAML itself, rather an implementation detail of the IdP.

Example SAML assertion:

```xml
<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion" ID="_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6" Version="2.0" IssueInstant="2014-07-17T01:01:48Z" Destination="https://sp.com/saml/acs" InResponseTo="C2dE3fH4iJ5kL6mN7oP8qR9sT0uV1w">
  <saml:Issuer>https://idp.com</saml:Issuer>
  <samlp:Status>
    <samlp:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"/>
  </samlp:Status>
  <saml:Assertion xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xs="http://www.w3.org/2001/XMLSchema" ID="pfx9e585951-b9dd-8c34-dbab-86cce4bd982b" Version="2.0" IssueInstant="2014-07-17T01:01:48Z">
    <saml:Issuer>https://idp.com</saml:Issuer><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
    <ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
  <ds:Reference URI="#pfx9e585951-b9dd-8c34-dbab-86cce4bd982b"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/><ds:DigestValue>KX1qZAXb9guMrtJBI0iTysdAZrI=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>Z30+fOmCKlWGR48uLkJgE7wjZIrf9oy87P9hw/RtaBSCsQGOnvXXaYHePfXQKVNkC3zPkeg935ue95CC5y4s7Ge7St0LkTAXMzF9h1vgB3xZCRtAo8zYOVvlh60Tcvuc/D5Odkbho8ZRKPymZteznrWxLnrss9dT8DvayFofCOU=</ds:SignatureValue>
<ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICajCCAdOgAwIBAgIBADANBgkqhkiG9w0BAQ0FADBSMQswCQYDVQQGEwJ1czETMBEGA1UECAwKQ2FsaWZvcm5pYTEVMBMGA1UECgwMT25lbG9naW4gSW5jMRcwFQYDVQQDDA5zcC5leGFtcGxlLmNvbTAeFw0xNDA3MTcxNDEyNTZaFw0xNTA3MTcxNDEyNTZaMFIxCzAJBgNVBAYTAnVzMRMwEQYDVQQIDApDYWxpZm9ybmlhMRUwEwYDVQQKDAxPbmVsb2dpbiBJbmMxFzAVBgNVBAMMDnNwLmV4YW1wbGUuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZx+ON4IUoIWxgukTb1tOiX3bMYzYQiwWPUNMp+Fq82xoNogso2bykZG0yiJm5o8zv/sd6pGouayMgkx/2FSOdc36T0jGbCHuRSbtia0PEzNIRtmViMrt3AeoWBidRXmZsxCNLwgIV6dn2WpuE5Az0bHgpZnQxTKFek0BMKU/d8wIDAQABo1AwTjAdBgNVHQ4EFgQUGHxYqZYyX7cTxKVODVgZwSTdCnwwHwYDVR0jBBgwFoAUGHxYqZYyX7cTxKVODVgZwSTdCnwwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQ0FAAOBgQByFOl+hMFICbd3DJfnp2Rgd/dqttsZG/tyhILWvErbio/DEe98mXpowhTkC04ENprOyXi7ZbUqiicF89uAGyt1oqgTUCD1VsLahqIcmrzgumNyTwLGWo17WDAa1/usDhetWAMhgzF/Cnf5ek0nK00m0YZGyc4LzgD0CROMASTWNg==</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature>
    <saml:Subject>
      <saml:NameID SPNameQualifier="http://sp.com" Format="urn:oasis:names:tc:SAML:2.0:nameid-format:transient">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>
      <saml:SubjectConfirmation Method="urn:oasis:names:tc:SAML:2.0:cm:bearer">
        <saml:SubjectConfirmationData NotOnOrAfter="2024-01-18T06:21:48Z" Recipient="http://sp.com" InResponseTo="C2dE3fH4iJ5kL6mN7oP8qR9sT0uV1w"/>
      </saml:SubjectConfirmation>
    </saml:Subject>
    <saml:Conditions NotBefore="2014-07-17T01:01:18Z" NotOnOrAfter="2024-01-18T06:21:48Z">
      <saml:AudienceRestriction>
        <saml:Audience>http://sp.com</saml:Audience>
      </saml:AudienceRestriction>
    </saml:Conditions>
    <saml:AuthnStatement AuthnInstant="2014-07-17T01:01:48Z" SessionNotOnOrAfter="2024-07-17T09:01:48Z" SessionIndex="C2dE3fH4iJ5kL6mN7oP8qR9sT0uV1w">
      <saml:AuthnContext>
        <saml:AuthnContextClassRef>freja-eid</saml:AuthnContextClassRef>
      </saml:AuthnContext>
    </saml:AuthnStatement>
    <saml:AttributeStatement>
      <saml:Attribute Name="userpersonalnumber" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
        <saml:AttributeValue xsi:type="xs:string">199601071489</saml:AttributeValue>
      </saml:Attribute>
      <saml:Attribute Name="mail" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic">
        <saml:AttributeValue xsi:type="xs:string">user@example.com</saml:AttributeValue>
      </saml:Attribute>
    </saml:AttributeStatement>
  </saml:Assertion>
</samlp:Response>
```

The assertion contains the following key elements:

- `saml:Issuer`: To make sure the SP can verify the response is from the expected IdP.
- `saml:AttributeStatement`: Contains user attributes like `userpersonalnumber` and `mail`. These attributes are used by the SP to validate the user exists and authorize access to resources.
- `saml:Subject`: Contains the `NameID` element, in this case **transient**. Which means it's a randomly generated identifier that is unique to the session. If it was **persistent**, it would be a stable identifier for the SP verify the user based on this string value.

This is all it takes for a successful SAML authentication flow. The SP receives the assertion, verifies the signature using the IdP's public key, and extracts the user's identity information. All using XML based protocols.
