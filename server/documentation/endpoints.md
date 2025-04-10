# Project Kallaxy Endpoints <!-- omit from toc -->

- [1. Users endpoints](#1-users-endpoints)
  - [1.1. POST /api/users -- User creation](#11-post-apiusers----user-creation)
  - [1.2. GET /api/users -- Get user info by ID (need an valid access token)](#12-get-apiusers----get-user-info-by-id-need-an-valid-access-token)
  - [1.3. PUT /api/users -- User info update](#13-put-apiusers----user-info-update)
  - [1.4. DELETE /api/users -- Delete all user's info (need a valid access token)](#14-delete-apiusers----delete-all-users-info-need-a-valid-access-token)
- [2. Authentification endpoints](#2-authentification-endpoints)
  - [2.1. POST /auth/login -- Get access token and refresh token](#21-post-authlogin----get-access-token-and-refresh-token)
  - [2.2. POST /auth/logout -- Logout a user](#22-post-authlogout----logout-a-user)
  - [2.3. POST /auth/refresh -- Refresh access token](#23-post-authrefresh----refresh-access-token)
  - [2.4. POST /auth/revoke -- Revoke a refresh token](#24-post-authrevoke----revoke-a-refresh-token)
  - [GET /auth/login -- Confirm user password](#get-authlogin----confirm-user-password)
- [3. Media endpoints](#3-media-endpoints)
  - [3.1. POST /api/media -- Create a new medium](#31-post-apimedia----create-a-new-medium)
  - [3.2. GET /api/media -- Get a medium's info by its title](#32-get-apimedia----get-a-mediums-info-by-its-title)
  - [3.3. GET /api/media/type -- Get all media based on given type](#33-get-apimediatype----get-all-media-based-on-given-type)
  - [3.4. PUT /api/media -- Update a medium's info](#34-put-apimedia----update-a-mediums-info)
  - [3.5. DELETE /api/media -- Delete a medium](#35-delete-apimedia----delete-a-medium)
- [4. Records endpoints](#4-records-endpoints)
  - [4.1. POST /api/records -- Create a new User-Medium Record](#41-post-apirecords----create-a-new-user-medium-record)
  - [4.2. GET /api/records -- Get all records by user's ID](#42-get-apirecords----get-all-records-by-users-id)
  - [GET /api/media\_records -- Get all user's record and related media](#get-apimedia_records----get-all-users-record-and-related-media)
  - [4.3. PUT /api/records -- Update a record's start and/ord end date](#43-put-apirecords----update-a-records-start-andord-end-date)
  - [4.4. DELETE /api/records -- Delete a record](#44-delete-apirecords----delete-a-record)
- [5. Password Reset endpoints (IN TEST MODE, NOT SECURE FOR PRODUCTION)](#5-password-reset-endpoints-in-test-mode-not-secure-for-production)
  - [5.1. POST /auth/password\_reset -- Step 1 : Ask for a reset token and reset link](#51-post-authpassword_reset----step-1--ask-for-a-reset-token-and-reset-link)
  - [5.2. GET /auth/password\_reset?token=xxxxxxxx -- Step 2 : Verify reset token](#52-get-authpassword_resettokenxxxxxxxx----step-2--verify-reset-token)
  - [5.3. PUT /auth/password\_reset -- Step 3 : Set a new password](#53-put-authpassword_reset----step-3--set-a-new-password)
- [External API endpoints (Server acts as a proxy)](#external-api-endpoints-server-acts-as-a-proxy)
  - [Books (on openLibrary.org)](#books-on-openlibraryorg)
    - [GET /external\_api/book/search -- Search for a book by title or by author](#get-external_apibooksearch----search-for-a-book-by-title-or-by-author)
    - [GET /external\_api/book/isbn](#get-external_apibookisbn)
    - [GET /external\_api/book/author](#get-external_apibookauthor)
    - [GET /external\_api/book/search\_isbn](#get-external_apibooksearch_isbn)
  - [Movies/Series](#moviesseries)
    - [GET /external\_api/movie\_tv/search\_movie](#get-external_apimovie_tvsearch_movie)
    - [GET /external\_api/movie\_tv/search\_tv](#get-external_apimovie_tvsearch_tv)
    - [GET /external\_api/movie\_tv/search](#get-external_apimovie_tvsearch)
    - [GET /external\_api/movie\_tv](#get-external_apimovie_tv)
  - [Videogames](#videogames)
    - [GET /external\_api/videogame/search](#get-external_apivideogamesearch)
    - [GET /external\_api/videogame](#get-external_apivideogame)
  - [Boardgames](#boardgames)
    - [GET /external\_api/boardgame/search](#get-external_apiboardgamesearch)
    - [GET /external\_api/boardgame](#get-external_apiboardgame)


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
> **WARNING : User's refresh tokens will be revoked and they need to login again to get new tokens.**

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
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 409 Conflict - username or email is already used by another user

-> *OK Response status code expected* : 

    200 OK
    /!\ If client receives 200 OK response, it also means that user's refresh token has been revoked and need to log in again for new refresh token /!\ 

-> *Response body example* :
>See resource [User](resources.md#user-resource)

### 1.4. DELETE /api/users -- Delete all user's info (need a valid access token)
-> *Description* :
>Delete logged user's row in server's database
>Empty response's body

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
### 2.1. POST /auth/login -- Get access token and refresh token
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
>See resource [User](resources.md#user-resource) and resource [Tokens](resources.md#tokens)

### 2.2. POST /auth/logout -- Logout a user
-> *Description* : 
> Logout a logged user by revoking all their refresh tokens
>Empty response's body

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
>Empty response's body

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

### GET /auth/login -- Confirm user password

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
* `title` - *string*
* `media_type` - *string*  
*Couple `title` & `media_type` needs to be unique across server's database*
* `creator` - *string*
* `pub_date` - *string*

> **OPTIONNAL**:
* `image_url` - *string*
* `metadata` - map[string]interface{} (according to media type see resources [metadata](resources.md#metadata-for-media))


*Example*:
```json
{
    "title": "The Fellowship of the Ring",
    "media_type": "book",
    "creator": "J.R.R Tolkien",
    "pub_date": "1954/02/01",
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - One to many required fields are missing in request's body
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 409 Conflict - A medium with the exact same title/media_type already exists in database

-> *OK Response status code expected* :

    201 Created

-> *Response body* :
> See resource [Medium](resources.md#media-resource)

### 3.2. GET /api/media -- Get a medium's info by its title
-> *Description* :
>Get info for a medium whose title and media_type is given in request body
>Respond with the found medium

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>**REQUIRED**:
* `title` - *string*
* `media_type` - *string*  

*Example*:
```json
{
    "title": "The Fellowship of the Ring",
    "media_type": "book"
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - Client didn't provide a correct query parameter.
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given title in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> See resource [Medium](resources.md#media-resource)

### 3.3. GET /api/media/type -- Get all media based on given type
-> *Description* :
>Get info for all media whose type is given in request body
>Respond with a list of media

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>**REQUIRED**:
* `media_type` - *string* 

*Example*:
```json
{
    "media_type": "book"
}
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
    "media": []Medium
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
* `pub_date` - *string*

> **OPTIONNAL**:
* `image_url` - *string*
* `metadata` - map[string]interface{} (according to media type see resources [metadata](resources.md#metadata-for-media))

>**`media_type` cannot be updated**  
>Even if a field is not updated, client still need to send old info (no comparison is done in server, all are replaced).

*Example*:
```json

{
    "medium_id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "title": "The Two Towers",
    "creator": "J.R.R Tolkien",
    "pub_date": "1954/02/01",
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
>Empty response's body

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
* `start_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#datetime))
* `end_date` - *string* (in format ISO 8601 datetime, see resource documentation [datetime](resources.md#datetime))

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

### GET /api/media_records -- Get all user's record and related media
-> *Description* :
> Find all records matching logged user (by user's id from access token) and all media related to those records
> Respond with a map[string][]MediumWithRecord

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
    "records": map[string][]MediumWithRecord
}
```
> See resource [MediumWithRecord](resources.md#media-with-record-resource)

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
>Empty response's body

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


## 5. Password Reset endpoints (IN TEST MODE, NOT SECURE FOR PRODUCTION)

### 5.1. POST /auth/password_reset -- Step 1 : Ask for a reset token and reset link
-> *Description* :
>Based on given user's email
* Server generates a unique, time-limited reset token (6h)
* Server stores token in database with user ID and expiration
* In production, server would emails the reset link to user
* In test, server responds with reset link

-> *Request headers* :
>None

-> *Request body* :
>**REQUIRED**:
* `email` - string

*Example*:
```json
{
     "email": "vincnt21@example.com"
}
```
-> *Error Response status code to handle* : 

    - None, client won't know if email exists in server's database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* **(IN TEST MODE ONLY)**:
```json
{
    "message": "Password reset initiated",
    "reset_link": "/auth/password_reset?token=vdsfsfe23456dfs",
    "token": "vdsfsfe23456dfs"
}
```

### 5.2. GET /auth/password_reset?token=xxxxxxxx -- Step 2 : Verify reset token
-> *Description* :
>Server verify if the token from query parameter exists, hasn't expired and hasn't already been used
> Respond with `valid` (*bool*) and `email` (*string*)

-> *Request headers* :
>REFRESH TOKEN

-> *Request body* :
>**REQUIRED**:
* `email` - string

*Example*:
```json
{
     "email": "vincnt21@example.com"
}
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Missing token in query parameters OR invalid / expired reset token

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
```json
{
    "valid": true,
    "email": "vincnt21@example.com"
}
```

### 5.3. PUT /auth/password_reset -- Step 3 : Set a new password
-> *Description* :
>New password is set for user (based on given reset token)
> All refresh token linked to user's ID will be revoked, user will need to login again to get new tokens.
>Respond with updated User

-> *Request headers* :
>None

-> *Request body* :
>**REQUIRED**:
* `token` - *string*
* `new_password` - *string*

*Example*:
```json
{
     "token": "KBTVMH4IAVEET7P6GIPUDKTYPS",
     "new_password": "qsdf5678"
}
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - invalid or expired reset token

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>See resource [User](resources.md#user-resource)


## External API endpoints (Server acts as a proxy)
### Books (on openLibrary.org)
#### GET /external_api/book/search -- Search for a book by title or by author
-> *Request query parameters:*  
> ?title=xxxx
> ?author=xxxxx

#### GET /external_api/book/isbn
-> *Request query parameters:*  
> ?isbn=xxxxx

#### GET /external_api/book/author
-> *Request query parameters:*  
> ?author=xxxxx

#### GET /external_api/book/search_isbn
-> *Request query parameters:*  
> ?key=xxxxx

### Movies/Series
#### GET /external_api/movie_tv/search_movie
-> *Request query parameters:*  
> ?query=xxxx

#### GET /external_api/movie_tv/search_tv
-> *Request query parameters:*  
> ?query=xxxx

#### GET /external_api/movie_tv/search
-> *Request query parameters:*  
> ?query=xxxx

#### GET /external_api/movie_tv
-> Request body:
movie_id string
tv_id string
language string

### Videogames
#### GET /external_api/videogame/search
> ?search=<title>&platforms=<platformsID>

#### GET /external_api/videogame
> ?id=xxxx

### Boardgames
#### GET /external_api/boardgame/search
> ?query=xxxx

#### GET /external_api/boardgame
> ?id=xxxx
