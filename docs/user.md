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

**Notes**
- already register to firebase, to get the idToken

## Get Firebase Claims
Key | Value
-- | --
Method | GET
URL | /me

Headers | Value
-- | --
Authorization | bearer {token}

Body | -
-- | --
No Body | -

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
photo_url | optional. string

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized

## Get Profile
Key | Value
-- | --
Method | GET
URL | /profile

Headers | Value
-- | --
Authorization | bearer {token}

Body | -
-- | --
No Body | -

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized

## Update Profile
Key | Value
-- | --
Method | PATCH
URL | /profile/:id

Note:
- id. user_id

Headers | Value
-- | --
Authorization | bearer {token}
Content-Type | application/json

Body | Description
-- | --
age | string
dob | date of birth. example: "2006-01-02T15:04:05Z"
sex | string. one of Male, Female
blood_type | string. one of A, B, O, AB
weight | float. example: 45.5
height | float. example: 155
activity_level | string. one of Sedentary, Lightly Active, Moderately Active, Very Active
allergies | string. example: "Allergies1,Allergies2"
ec_relation | string. one of Husband, Wife, Mother, Father etc.
ec_name | string
ec_country_code | example: 62, 65, etc
ec_phone | numeric. 9-12 charactes. not start with 0. example: 81212341234

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized

## Delete Profile
Key | Value
-- | --
Method | DELETE
URL | /profile/:id

Note:
- id. user_id

Headers | Value
-- | --
Authorization | bearer {token}
Content-Type | application/json

Body | -
-- | --
No Body | -

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

Headers | -
-- | --
No Header | -

Body | Description
-- | --
email | string

Status Code | Value
-- | --
200 | Success
400 | Bad Request

## Update Password
Key | Value
-- | --
Method | POST
URL | /credentials/update-password

Headers | -
-- | --
No Header | -

Body | Description
-- | --
user_id | string
reset_token | string
password | string

Status Code | Value
-- | --
200 | Success
400 | Bad Request

## Upload Photo Profile
Key | Value
-- | --
Method | PATCH
URL | /profile/:id/photo

Note:
- id. user_id

Headers | Value
-- | --
Authorization | bearer {token}
Content-Type | multipart/form-data

Body | Description
-- | --
file | photo to be uploaded

Response | Description
-- | --
data | photo url

Status Code | Value
-- | --
200 | Success
400 | Bad Request
401 | Unauthorized