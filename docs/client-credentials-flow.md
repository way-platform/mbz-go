# OAuth Authentication (Client Credentials Flow)

Client Credentials

Our APIs use standard OAuth 2.0 to allow authorization of personal data. Read on if you want to know how OAuth 2.0 works and how to integrate your application with our security components. You can find the detailed specification for the OAuth 2.0 Authorization Framework [here](https://tools.ietf.org/html/rfc6749).

**Note:** This document describes the client credentials code flow. Some of our products use the authorization code flow, which is documented [here](https://developer.mercedes-benz.com/content-page/oauth-documentation).

## Client Credentials

Your client credentials carry many privileges, so be sure to keep them secure! Therefore, make sure that your application fulfills the following prerequisites:

- It must keep the client credentials on a safe place in the backend (on the server).

- It must not expose any of the client credentials (client id, client secret, access tokens) to the client.


## Client Credentials Flow

If an API supports the Client Credentials Flow, The client_credentials grant type is used by clients to obtain an access token outside of the context of a user. This is typically used by clients to access resources about themselves rather than to access a user’s resources.

These are the steps that the flow executes:

1. Request to get an access token

2. Use the access token to call the API


### From the developer’s perspective

In the following, we demonstrate how the client credentials code flow looks from a developer’s perspective who is developing an application using our provided APIs.

#### 1. Request to get an access token

In order to initiate the end user’s authorization, you should call our access token endpoint. The client must authenticate using the HTTP Basic method and provide the clientId and the clientSecret (_<insert_your_client_id_here>:<insert_your_client_secret_here>_) encoded with BASE64 in the HTTP Authorization header.

The following parameters are used in the **--data** value:

- **grant_type=client_credentials:** This indicates for the token endpoint to use the OAuth 2.0 Client Credentials Flow for this request.

- **scope=openid**: This scope is needed in order to receive a valid token.


You will then receive the OAuth access token in the server response. Note that the **expires_in** field in the response represents the validity period of the access token in seconds. Per default it is 3600s.

Example Request

```http
curl --request POST 'https://ssoalpha.dvb.corpinter.net/v1/token' \
--header 'Authorization: Basic <insert_your_base64_encoded_client_id_and_client_secret_here>' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data 'grant_type=client_credentials' \
--data 'scope=openid'
```

Example Authorization Header

```http
Authorization: Basic VGhpc0lzWW91ckNsaWVudElkOlRoaXNJc1lvdXJDbGllbnRTZWNyZXQ=
```

Example Response

```json
{
   "access_token":"<your_access_token>",
   "token_type":"bearer",
   "expires_in":3599,
   "id_token":"<your_id_token>"
}
```

#### 2. Use the access token to call the API

Now you can use the access token to call the API as long as it is not expired. Add the provided token to the authorization header to your API request.

Typical error cases:

|                     |                                                |
| ------------------- | ---------------------------------------------- |
| **HTTP Error Code** | **Description**                                |
| 401                 | The given access token is not valid (anymore). |
