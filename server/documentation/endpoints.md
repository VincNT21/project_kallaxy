# Project Kallaxy Endpoints <!-- omit from toc -->


## 1. Users endpoints

### 1.1. POST /api/users -- User creation
-> *Description* : 
>Create a new user in **users** table
>Respond with a User struct

-> *Request body* :
> **REQUIRED**:
* unique `username` - *string*
* `password` - *string* 
* unique `email` - *string*

*Example*:
```json
{
    "username": "VincNT21",
    "password": "12345abcde",
    "email": "vincnt21@example.com"   
}

```
-> *Error Response status codes to handle* : 

    - 400 Bad Request : one or many fields are missing in request
    - 409 Conflict : username or email is already used by another user

-> *OK Response status code expected* : 

    201 Created

-> *Response body example* :
>See resource [User](resources.md#user-resource)

### 1.2. GET /api/users -- Get user info by ID (need an valid access token)
-> *Description* :
>Get all info from database about logged user  
>Respond with a User struct

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>See resource [User](resources.md#user-resource)

### 1.3. PUT /api/users -- User info update
-> *Description* : 
> Update username/password/email for a logged-in user.  
> Respond with a User struct
> **WARNING : User's refresh tokens will be revoked and they need to login again.**

-> *Request headers* : 
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>**REQUIRED**:  
* `username` - string
* `password` - string
* `email` - string
>Even if a field is not updated, client still need to send old info (no comparison is done in server, all are replaced).

*Example*:
```json
{
    "username": "VincNT21",
    "password": "12345ghjk",
    "email": "vincnt21@example.com" 
}

```
-> *Error Response status code to handle* : 

    - 400 Bad Request - One to many fields are missing in request
    - 401 Unauthorized - Means that access token is expired, client should fetch **POST /auth/refresh** to get a new access token

-> *OK Response status code expected* : 

    200 OK
    /!\ If client receives 200 OK response, it also means that user's refresh token has been revoked and need to log in again for new refresh token /!\ 

-> *Response body example* :
>See resource [User](resources.md#user-resource)

### 1.4. DELETE /api/users -- Delete all user's info (need a valid access token)
-> *Description* :
>Delete logged user's row in server's database
>Response body empty

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>None


## 2. Authentification endpoints
### 2.1. POST /auth/login -- Authentification
-> *Description* : 
> Login user by checking given email/password, create Refresh Token (valid for 60 days) stored in server's database and a Access Token (valid for 1 hour) not stored. 
> Respond with both tokens and the logged user's info.

-> *Request body* :
>**REQUIRED**:
* valid `username` *string*
* valid `password` *string*

*Example*:
```json
{
    "username": "VincNT21",
    "password": "12345abcde"
}

```

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Given username/password does not match.

-> *OK Response status code expected* : 

    201 Created

-> *Response body example* :
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "username": "VincNT21",
    "email": "vincnt21@example.com",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrYWxsYXh5Iiwic3ViIjoiODFjMWNiMGQtYmJkYi00ZmFhLWFlZGUtYmQzNzFhNGFiNzIyIiwiZXhwIjoxNzQzMDIxMTMyLCJpYXQiOjE3NDMwMTc1MzJ9.1PUE_93e6pXaLwjZiMIfr5DAcxTxE4jEIiRftQuJptI",
    "refresh_token": "176ddabd5f4c932b8cda583e00b620a05242187680002a071e8a13c4e2b0b14"
}
```

### 2.2. POST /auth/logout -- Logout a user
-> *Description* : 
> Logout a logged user by revoking all their refresh tokens
>Response body empty

-> *Request headers* : 
> A valid access token in "Authorization" header
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No refresh token associated to user's ID (ID comes from access token) found in database

-> *OK Response status code expected* : 

    204 No Content

### 2.3. POST /auth/refresh -- Refresh access token
-> *Description* : 
>If given Refresh Token is still valid and not revoked, create a new Access Token and a new Refresh Token.
>Given refresh token will be revoked
> Respond with both tokens

-> *Request headers* : 
>A valid refresh token (string) in "Authorization" header
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Given refresh token doesn't exist in server's database or has been revoked or has expired. Client should fetch **POST /auth/login** to get a new refresh token.

-> *OK Response status code expected* : 

    201 Created

-> *Response body example* :
```json
{
    "access_token": "<access_token>",
    "refresh_token": "<refresh_token>"
}
```
> See resource [Tokens](resources.md#tokens) for more details

### 2.4. POST /auth/revoke -- Revoke a refresh token
-> *Description* : 
>Revoke a refresh token in server's database
>Response body empty

-> *Request headers* : 
>A valid refresh token (string) in "Authorization" header
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - There is a problem with "Authorization" header: missing or malformed
    - 404 Not Found - Given refresh token doesn't exist in server's database

-> *OK Response status code expected* : 

    204 No Content

-> *Response body example* :
>None


## 3. Media endpoints

### 3.1. POST /api/media -- Create a new medium
-> *Description* :
>Create a new medium in server's database
>Respond with the created medium

-> *Request headers* :
> A valid Bearer access token in "Authorization" header.  
> See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>**REQUIRED**:
* unique `title` - *string*
* `media_type` - *string*
* `creator` - *string*
* `release year` - *int32*

> **OPTIONNAL**:
* `image_url` - *string*
* `metadata` - json struct (according to media type see resources [metadata](resources.md#metadata-for-media))


*Example*:
```json
{
    "title": "The Fellowship of the Ring",
    "media_type": "book",
    "creator": "J.R.R Tolkien",
    "release_year": 1954,
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - One to many required fields are missing in request's body
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 409 Conflict - A medium with the exact same title already exists in database

-> *OK Response status code expected* :

    201 Created

-> *Response body* :
> See resource [Medium](resources.md#media-resource)

### 3.2. GET /api/media?title=*xxx* -- Get a medium's info by its title
-> *Description* :
>Get info for a medium whose title is given in request query parameters
>Respond with the found medium

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request query parameters* :
>"?title=<medium_title>"  
Note that medium title is case insensitive (lowered before database query) BUT spaces (encoded with + or %20) and special characters matters for server search

*Example*:
```
/api/media?title=The+Fellowship+Of+The+Ring
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Client didn't provide a correct query parameter.
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given title in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> See resource [Medium](resources.md#media-resource)

### 3.3. GET /api/media?type=*xxx* -- Get all media based on given type
-> *Description* :
>Get info for all media whose type is given in request query parameters
>Respond with a list of media

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request query parameters* :
>"?type=<media_type>"  
Note that media type is case insensitive (lowered before database query) BUT spaces (encoded with + or %20) and special characters matters for server search

*Example*:
```
/api/media?type=book
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Client didn't provide a correct query parameter.
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No media of given type in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
```json
{
    "media": []medium
}
```
> See resource [Medium](resources.md#media-resource)

### 3.4. PUT /api/media -- Update a medium's info
-> *Description* :
> Change some info about a specified medium (by medium's id)  
> Respond with updated medium

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> **REQUIRED**:
* `medium_id` - *string* (in format UUIDv4, see [resource documentation](resources.md#uuid))
* unique `title` - *string*
* `creator` - *string*
* `release year` - *int32*

> **OPTIONNAL**:
* `image_url` - *string*
* `metadata` - json struct (according to media type see resources [metadata](resources.md#metadata-for-media))

>**`media_type` cannot be updated**  
>Even if a field is not updated, client still need to send old info (no comparison is done in server, all are replaced).

*Example*:
```json

{
    "medium_id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "title": "The Two Towers",
    "creator": "J.R.R Tolkien",
    "release_year": 1954,
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - Medium_id not in UUIDv4 format
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given ID found in database
    - 409 Conflict - A medium with the exact same title already exists in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> See resource [Medium](resources.md#media-resource)

### 3.5. DELETE /api/media -- Delete a medium
-> *Description* :
>Delete a medium's info in database, based on given medium's ID
>Response body empty

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> **REQUIRED**
* `medium_id` - *string* (in format UUIDv4, see [resource documentation](resources.md#uuid))

*Example*:
```json
{
    "medium_id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc"
}
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Medium_id not in UUIDv4 format
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given ID found in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>None

## 4. Records endpoints

### 4.1. POST /api/records -- Create a new User-Medium Record
-> *Description* :
> Create a record linking a user (based on user's id from access token) and a medium (based on given medium's id from body request).  
> Respond with this newly created record

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> **REQUIRED**: 
* `medium_id` - *string* (in format UUIDv4, see resource documentation [UUID](resources.md#uuid))   
> **OPTIONNAL**: 
* `start_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#iso-8601-datetime))
* `end_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#iso-8601-datetime))

*Example*:
```json
{
    "medium_id": "3b75af06-e596-42ce-a953-bf235dfc9102",
    "start_date": "2025-03-26T14:20:23.525332",
    "end_date": "2025-03-31T08:47:29.205805",
}
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Request's body missing medium_id OR medium_id not in UUIDv4 format OR request's dates not in ISO 8601 format OR request's start date is before request's end date
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No user or medium found in database with given ID
    - 409 Conflict - A record with the same user-medium couple already exists in database

-> *OK Response status code expected* :

    201 Created

-> *OK Response body example* :
>See resource [Record](resources.md#record-resource)


### 4.2. GET /api/records -- Get all records by user's ID
-> *Description* :
> Find all records matching logged user (by user's id from access token)
> Respond with a list of those records

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No record found for logged user

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
```json
{
    "records": []Record
}
```
> See resource [Record](resources.md#record-resource)

### 4.3. PUT /api/records -- Update a record's start and/ord end date 
-> *Description* :
> Modify a already-existing record's start date and/or end date (based on given record's ID)
> Respond with the updated record

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
> *OPTIONNAL*:
* `start_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#iso-8601-datetime))
* `end_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#iso-8601-datetime))

*Example*:
```json
{
    "end_date": "2025-03-31T08:47:29.205805"
}
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Start date (given or already existing) is before end date (given or already existing)
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No record found with given record's ID

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> See resource [Record](resources.md#record-resource)

### 4.4. DELETE /api/records -- Delete a record
-> *Description* :
>Delete a record based on given record's ID  
>No response

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>**REQUIRED**:
* `record_id` - *string* (in format UUIDv4, see resource documentation [UUID](resources.md#uuid))  

*Example*:
```json
{
    "record_id": "3b75af06-e596-42ce-a953-bf235dfc9102"
}
```
-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No record found with given record's ID

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>Empty