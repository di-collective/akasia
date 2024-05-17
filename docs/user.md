# User Service

## Register/Login
Key | Value
-- | --
Method | POST
URL | /credentials/firebase-auth

Query Params | Value
-- | --
idToken | {idToken}

Headers | Value
-- | --
Content-Type | application/json

Body | Description
-- | --
provider | {provider}
email | {email}
password | {password}
repeat_password | {repeat_password}

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized

## Create Profile
Key | Value
-- | --
Method | POST
URL | /profile

Headers | Value
-- | --
Authorization | bearer {token}
Content-Type | application/json

Body | Description
-- | --
name | string
country_code | example: 62, 65, etc
phone | numeric. 9-12 charactes. not start with 0. example: 81212341234
nik | optional. numeric. 16 characters

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized

## Forgot Password
Key | Value
-- | --
Method | POST
URL | /credentials/forgot-password

**Headers**

No Header

Body | Description
-- | --
email | string

Status Code | Value
-- | --
200 | Success
400 | Bad Request
