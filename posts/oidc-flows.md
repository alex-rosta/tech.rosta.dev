---
title: "OpenID Connect Flows"
created: 2025-07-24
updated: 2025-07-24
tags: "oidc, authentication, identity"
---
Let's go over some OIDC flows, what makes them different, and when to use each one.
![image](https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/OpenID_logo_2.svg/1200px-OpenID_logo_2.svg.png)

## Basics

OpenID Connect (OIDC) is an authentication layer built on top of OAuth 2.0. It allows clients to verify the identity of users based on the authentication performed by an authorization server. OIDC introduces the concept of an ID token, which is a JSON Web Token (JWT) that contains user information.

To understand OIDC on a basic level, we can break a `.well-known/openid-configuration` document together, which describes the endpoints and capabilities of an OIDC provider.

### [Google's `.well-known/openid-configuration`](https://accounts.google.com/.well-known/openid-configuration)

```json
{
"issuer": "https://accounts.google.com",
"authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
"device_authorization_endpoint": "https://oauth2.googleapis.com/device/code",
"token_endpoint": "https://oauth2.googleapis.com/token",
"userinfo_endpoint": "https://openidconnect.googleapis.com/v1/userinfo",
"revocation_endpoint": "https://oauth2.googleapis.com/revoke",
"jwks_uri": "https://www.googleapis.com/oauth2/v3/certs",
"response_types_supported": [
"code",
"token",
"id_token",
"code token",
"code id_token",
"token id_token",
"code token id_token",
"none"
],
"subject_types_supported": [
"public"
],
"id_token_signing_alg_values_supported": [
"RS256"
],
"scopes_supported": [
"openid",
"email",
"profile"
],
"token_endpoint_auth_methods_supported": [
"client_secret_post",
"client_secret_basic"
],
"claims_supported": [
"aud",
"email",
"email_verified",
"exp",
"family_name",
"given_name",
"iat",
"iss",
"name",
"picture",
"sub"
],
"code_challenge_methods_supported": [
"plain",
"S256"
],
"grant_types_supported": [
"authorization_code",
"refresh_token",
"urn:ietf:params:oauth:grant-type:device_code",
"urn:ietf:params:oauth:grant-type:jwt-bearer"
]
}
```

## Authorization Code Flow

The most common flow, used for server-side applications that are capable of secure storage of client secrets.
Parameters additional to the `.well-known/openid-configuration`:

- `Client ID` not sensitive, identifies the client application.
- `Client Secret` secretly shared to the client application. **Should be kept confidential at all times.**
- `Redirect URI` where the authorization server will send the user after authentication. will have to be whitelisted in the OIDC provider.

1. User accesses a client application, which requires authentication.
2. If the app is properly written it will already have cached and parsed the `.well-known/openid-configuration` document.
3. The app redirects the user to the `authorization_endpoint` with the following parameters:
   Example: `https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=CLIENT_ID&redirect_uri=https://yourapp.com/callback&scope=openid`
4. After authentication and consent from user, the authorization server redirects back to the `redirect_uri` with an authorization code.
   Example: `https://yourapp.com/callback?code=AUTHORIZATION_CODE`
5. Now the sensitive part happens, the client application exchanges the authorization code for an access token and ID token by making a POST request to the `token_endpoint` with the following parameters:
   - `grant_type=authorization_code`
   - `code=AUTHORIZATION_CODE`
   - `redirect_uri=https://yourapp.com/callback`
   - `client_id=CLIENT_ID`
   - `client_secret=CLIENT_SECRET`
   Obviously, this request should **always** happen server-side, never in the browser. `client_secret` should never be exposed to the client.
6. The authorization server responds with an access token and ID token.
   Example response:

   ```json
   {
     "access_token": "ACCESS_TOKEN",
     "id_token": "ID_TOKEN_JWT",
     "expires_in": 3600,
     "token_type": "Bearer"
   }
   ```

   The application should now verify the signature of `ID_TOKEN_JWT` using the public keys from the `jwks_uri` endpoint, and then decode it to get user information.
   If we decode the `ID_TOKEN_JWT`, we can see the user information:

   ```json
   {
     "iss": "https://accounts.google.com",
     "sub": "USER_ID",
     "email": "USER_EMAIL",
     "exp": 1712345678,
     "iat": 1712342078
   }
   ```

### What if an application is not capable of securely storing the client secret?

In the case of SPA apps, mobile etc. The application should use the **PKCE (Proof Key for Code Exchange)** extension to OIDC. PKCE allows public clients to securely authenticate without needing a client secret.
The flow is similar to the Authorization Code Flow, but with the subtraction of `client_secret` and addition of a `code_challenge` and `code_verifier`.

1. The client generates a `code_verifier`, which is a random string each time.
2. The client then creates a `code_challenge` by hashing the `code_verifier` using SHA-256.
3. The client sends the `code_challenge` along with the authorization request to the authorization server.
   Example: `https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=CLIENT_ID&redirect_uri=https://yourapp.com/callback&scope=openid&code_challenge=CODE_CHALLENGE&code_challenge_method=S256` (make sure this is supported by the OIDC provider, which is the case for Google).
4. After the user authenticates and the authorization server redirects back to the `redirect_uri` with an authorization code, the client sends a POST request to the `token_endpoint` with the following parameters:
   - `grant_type=authorization_code`
   - `code=AUTHORIZATION_CODE`
   - `redirect_uri=https://yourapp.com/callback`
   - `client_id=CLIENT_ID`
   - `code_verifier=CODE_VERIFIER`
5. This is the crucial part: the authorization server verifies that the `code_verifier` matches the `code_challenge` sent in the initial request. If they match, it responds with an access token and ID token. This ensures that the request is coming from the same client that initiated the flow, preventing interception attacks.

## Try to always stay OIDC conformant

OIDC is a standard, and while it allows for some flexibility, it's best to stick to the standard flows and parameters. This ensures compatibility with various OIDC providers and libraries. If you find yourself needing to deviate from the standard, consider whether it's worth rewriting the whole auth the next time another OIDC provider is knocking on the door.
