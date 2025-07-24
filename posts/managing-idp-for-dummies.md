---
title: "IdP handbook/rulebook for dummies"
created: 2025-07-24
updated: 2025-07-24
tags: "idp, authentication, identity"
---

## Introduction

Managing an Identity Provider (IdP) can and will be complex. This guide aims to simplify the process, providing pain points, best practices and common pitfalls to avoid.
Having been responsible for managing an IdP in a 15k+ user environment with 70+ applications, I've learned a lot about what approaches work and which don't.

![image](https://www.onelogin.com/blog/wp-content/uploads/2023/05/advanced-authentication-blog-image.jpg.optimal.jpg)

## Authentication Methods

### Third Party Authentication

*The best* by far approach, no user onboarding, no password resets.
If possible, use national eID solutions, such as [Freja eID](https://frejaeid.atlassian.net/wiki/spaces/DOC/pages/2162802/Authentication+Service) or [eIDAS for EU](https://www.swedenconnect.se/). (Sorry only swedish docs for eIDAS, refer to your country's documentation)

### Provide Multiple Authentication Methods

This is a good idea, but find a common claim these providers return from their API and make sure it's always present in your catalog.
If you use Freja and BankID, make sure Personnummer is stored in your catalog, and make that the common identifier if looking up users before populating claims.

## Claims

### Enriched Claims or Pass-Through Claims

If you rely on external authentication methods, such as [Siths eID](https://docs.grandid.com/SITHSeID), [Freja](https://frejaeid.atlassian.net/wiki/spaces/DOC/pages/2162802/Authentication+Service) or any other, you will have been given a set of initial claims.

Now, should you verify these claims against your own user catalog, or should you pass them through as they are?

In my opinion, if this question is asked even once you should have two instances of your IdP, one for enriched claims and one for pass-through claims.

- **Enriched Claims**: If you *can* guarantee the user is always present in your cataloges. Internal employees etc.
- **Pass-Through Claims**: If you *cannot* guarantee the user is always present. Guest users, external partners, etc.

### Standardization

Adhere to a standard set of claims. Take inspiration from tech giants in your field, how do they solve this problem? Is there a standard in your country or industry?
Is your org big enough to write your own standard?

- [Sambi-standard](https://wiki.federationer.internetstiftelsen.se/pages/viewpage.action?pageId=46465316)
- [Microsoft-standard](https://learn.microsoft.com/en-us/entra/identity-platform/reference-saml-tokens)

This is the most important thing for scalability.

The best way I have found to do this is to set a standard profile for all applications, and then allow deviation if necessary.

Find attributes that are common across your organization. Such as email, an employee ID etc. Build group-based claims around departments, roles etc. This is more IAM related, but the bottleneck is always going to be how much claims varies across applications.
If you can adhere to a standard early on, and build authorization groups around already existing organzational structures, you are bound for success.

### Deviation from Standard

This will always happen right. Some poorly written *absolutely* critical HR system will require a custom claim not present in your catalog, or maybe even now the IdP needs to support guest users not even present in your foundational directory.

Keep them documented, and make sure to separate them from the standard profile.
If you deviate too many times from standard, it will be a nightmare to maintain.

## Metadata Parsing

### URL based metadata

Always use this, SP's should cache SAML metadata or OIDC configuration documents, and you should be able to update them without having to reconfigure the SP.
Imagine rotating your signing key, and having to update hundreds of SP's IdP metadata files...

## Third Party Lookups

## Automation

## Authorization
