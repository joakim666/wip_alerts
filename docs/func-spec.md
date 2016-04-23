# Ismonitor / Isalerts / Ismoalerts / Ismolerts

## Registration

When the app is started check for a {account id, access token and refresh token} in local secure storage

### If no info is present

1. POST /accounts { device_id: <uuid of device>, device_type: '<ios|android|etc>', device_info: {<all relevant info like model, os version etc>}
    => get account_id (uuid) back
2. POST /token {grant_type: 'account', account_id: <account id>}
    => get refresh and access token back
3. Save account id, access token and refresh token in local secure storage

### If info is present

1. GET /ping
    => HTTP status '204 No Content' if access token is still valid otherwise renew access token as below
    
#### Renew access token

1. POST /renewals { account_id: <account_id uuid>, device_info: {<all relevant info like model, os version etc>}
    => get renewal_id (uuid) back
2. POST /token {grant_type: 'renewal', renewal_id: <renewal id>
    => get access token back
3. Use access token for all future request

## Authentication and Authorization

Authentication is done by using JWT tokens with JWE encrypted content. This way authentication and authorization can be done without having to touch the database.

There are a two different kinds of tokens: Refresh tokens and access tokens, a la OAuth 2. Refresh tokens can only be used to retrieve an access token which has a shorter time to live than the refresh token. The access token can then be used to access the rest of the API.

Tokens can have two different roles: 'user' and 'publisher'. The user role is used when using the app. The publisher role is used by reporting applications that only publish data but don't consume it.

The tokens also have one or more capabilities describing what the user of the token is allowed to do. This capabilities are in addition to the capabilities implied by the role.

## Refresh token vs access token

It is done this way to limit the checking against the token revocation list. Access tokens are not checked against the revocation list and will granted access during their time to live period. Refresh tokens on the other hand are checked, so when the access token has expired and the client request a new access token using the refresh token, the refresh token is checked against the revocation list.

The expiration time for the access tokens will be set serverside and will be something like 30 minutes.

## Setting up a publishing application

### In the App

1. Create a new refresh token with the publish-role (by using the api-keys endpoint)
2. Copy the created refresh token and paste it into the reporting program configuration file

### The reporting program

1. When starting uses the refresh token from the configuration file to request an access token
2. The access token is when used for all subsequent calls to the API until the access token expires. Then a new one is requested using the request token from the configuration.


### Refresh token example

{
  iat: 1416929061, // when the token was issued (seconds since epoch)
  jti: "802057ff9b5b4eb7fbb8856b6eb2cc5b", // a unique id for this token (for revocation purposes)
  sub: "<uuid>", // the unique uuid identifying the user
  type: "refresh_token", // the type of the token, 'refresh_token' or 'access_token'
  scopes: {
    roles: ['user'], // what roles this token has
    capabilities: ['access_token'], // what capabilities this token has
  }
}

### Access token example

{
  iat: 1416929061, // when the token was issued (seconds since epoch)
  jti: "802057ff9b5b4eb7fbb8856b6eb2cc5b", // a unique id for this token used for audit
  sub: "<uuid>", // the unique uuid identifying the user
  type: "access_token", // the type of the token, 'refresh_token' or 'access_token'
  scopes: {
    roles: ['user'], // what roles this token has
    capabilities: [''], // what capabilities this token has
  }
}



------

# API-endpoints

## POST /accounts

## POST /token

## GET /ping

## POST /renewals

## 

------

https://apiblueprint.org/

https://github.com/apiaryio/dredd


-----

API Blueprint below:

FORMAT: 1A

# Ismolerts

The Ismolerts API is used both by the iOS application and all reporting applications.

# Group Authentication

Resources related to authentication and token handling.

## Accounts resource [/accounts]

## Create a new account [POST]

This is used to create a new account, i.e. the first time the app is started on the phone.

+ device_id (string)    - uuid of device
+ device_type (string)  - ios|android|etc
+ device_info (map)     - all relevant info like model, os version etc

+ Request (application/json)

    {
        "device_id": "<uuid>",
        "device_type": "ios",
        "device_info": {
            "model": "iPhone6s",
            "os": "9.3"   
        }
    }

+ Response 201 (application/json)

    {
        "account_id": "<account uuid>"
    }

## Renewal resource [/renewals]

## Request a renewal [POST]

Request a renewal id that can be used to get a new access token.

+ account_id (string)       - the account id
+ refresh_token (string)    - the refresh token
+ device_type (string)      - ios|android|etc
+ device_info (map)         - all relevant info like model, os version etc

+ Request (application/json)

    {
        "account_id": "<uuid>",
        "refresh_token": "refresh_token",
        "device_type": "ios",
        "device_info": {
            "model": "iPhone6s",
            "os": "9.3"   
        }
    }

+ Response 201 (application/json)

    {
        "renewal_id", "<renewal uuid>"
    }

## Token resource [/token]

### Request a token [POST]

This is used to create a new token.

+ grant_type (string)               - the grant type should be 'renewal' or 'account'
+ account_id (optional,string)      - the account id. Required when grant_type is 'account'
+ renewal_id (optional,string)      - the renewal id. Required when grant_type is 'renewal'

+ Request (application/json)

    {
        "grand_type": "renewal",
        "renewal_id": "wkalj324jsdfkjl3242"
    }

+ Response 201 (application/json)

    {
        {
            iat: 1416929061, // when the token was issued (seconds since epoch)
            jti: "802057ff9b5b4eb7fbb8856b6eb2cc5b", // a unique id for this token (for revocation purposes)
            sub: "<uuid>", // the unique uuid identifying the user
            type: "refresh_token", // the type of the token, 'refresh_token' or 'access_token'
            scopes: {
                roles: ['user'], // what roles this token has
                capabilities: ['access_token'], // what capabilities this token has
            }
        }
    }

## Ping resource [/ping]

### Ping the service [GET]

Ping the service to make sure it's up and to validate the access token.

+ Response 204 (application/json)

If the access token is valid.


