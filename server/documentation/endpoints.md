# Project Kallaxy Endpoints <!-- omit from toc -->


## 1. Users Endpoints

### 1.1. POST /api/users -- User creation
-> *Description* : 
>Create a new user in **users** table

-> *Request body* :
>A unique username (string), a password (string) and a unique email (string)

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
>Get and respond all info from database about logged user

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
> Update username/password/email for a logged-in user. **WARNING : User's refresh tokens will be revoked and they need to login again.**

-> *Request headers* : 
> A valid access token in "Authorization" header

*Example*:
```json
{
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiZjJjMmVlMWQtYWExZS00YzBiLTliNmEtODcyMmY5OWE1ZWQwIiwiZXhwIjoxNzQxNjIxODYyLCJpYXQiOjE3NDE2MTgyNjJ9._9-QuSMwwy8zEAgWyq7gcayyRUzN-DDXolWz8VmXIMc"
}

``` 

-> *Request body* :
>A username (string), a password (string) and a email (string). 
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
> Login user by checking given email/password, create Refresh Token (valid for 60 days) stored in server's database and a Access Token (valid for 1 hour) not stored. Both tokens are sent back to client, along with the logged user's info.

-> *Request body* :
> user's username (string) and user's password (string)

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
    "access_token": "<access_token>",
    "refresh_token": "<refresh_token>"
}
```

### 2.2. POST /auth/logout -- Logout a user
-> *Description* : 
> Logout a logged user by revoking all their refresh tokens

-> *Request headers* : 
> A valid access token in "Authorization" header

*Example*:
```json
{
    "Authorization": "Bearer <access_token>"
}

``` 

-> *Request body* :
>None

-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No refresh token associated to user's ID (ID comes from access token) found in database

-> *OK Response status code expected* : 

    204 No Content

### 2.3. POST /auth/refresh -- Refresh access token
-> *Description* : 
>If given Refresh Token is still valid and not revoked, create a new Access Token and a new Refresh Token. Both tokens are sent back to client.

-> *Request headers* : 
>A valid refresh token (string) in "Authorization" header

*Example*:
```json
{
      "Authorization": "Bearer <refresh_token>"
}

``` 

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

### 2.4. POST /auth/revoke -- Revoke a refresh token
-> *Description* : 
>Revoke a refresh token in server's database

-> *Request headers* : 
>A valid refresh token (string) in "Authorization" header

*Example*:
```json
{
      "Authorization": "Bearer <refresh_token>"
}

``` 

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

-> *Request headers* :
> A valid Bearer access token in "Authorization" header.  
> See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>Title, type, creator (string), release year (int32), image URL (string), some metadata according to media type (see resources documentation)
All fields required, "type" field should be unique across database.

*Example*:
```json
{
    "title": "The Fellowship of the Ring",
    "type": "book",
    "creator": "J.R.R Tolkien",
    "release_year": "1954",
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - A medium with the exact same title already exists in database
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token

-> *OK Response status code expected* :

    201 Created

-> *Response body* :
> Returning the created medium
> See resource [Medium](resources.md#media-resource)

### 3.2. GET /api/media?title=*xxx* -- Get a medium's info by its title
-> *Description* :
>Get info for a medium whose title is given in request query parameters

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request query parameters* :
>"?title=<medium_title>"  
Note that medium title is case insensitive (lowered before database query) BUT spaces (encoded with %20 or +) and special characters matters for server search

*Example*:
```
/api/media?title=Fellowship%20Of%20The%20Ring
```
-> *Error Response status code to handle* : 

    - 400 Bad Request - Client didn't provide a correct query parameter.
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given title in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> Returning the searched medium
> See resource [Medium](resources.md#media-resource)

### 3.3. GET /api/media?type=*xxx* -- Get all media based on given type
-> *Description* :
>Get info for all media whose type is given in request query parameters

-> *Request headers* :
>A valid Bearer access token in "Authorization" header  
>See resource [Authorization header](resources.md#authorization-header)

-> *Request query parameters* :
>"?type=<media_type>"  
Note that media type is case insensitive (lowered before database query) BUT spaces (encoded with %20 or +) and special characters matters for server search

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
> A json list of medium
> See resource [Medium](resources.md#media-resource)

### 3.4. PUT /api/media -- Update a medium's info
-> *Description* :
> Change some info about a specified medium (by medium's id). 

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>First, an "id" field with medium's id.
>Then, same fields as for medium creation (see above). 
>**except for type which cannot be updated**.
>Even if a field is not updated, client still need to send old info (no comparison is done in server, all are replaced).

*Example*:
```json

{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "title": "The Fellowship of the Ring",
    "creator": "J.R.R Tolkien",
    "release_year": "1954",
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

-> *Error Response status code to handle* : 

    - 400 Bad Request - A medium with the exact same title already exists in database
    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given ID found in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
> Returning updated medium
> See resource [Medium](resources.md#media-resource)

### 3.5. DELETE /api/media -- Delete a medium
-> *Description* :
>Delete a medium's info in database, based on given medium's ID

-> *Request headers* :
>A valid Bearer access token in "Authorization" header 
>See resource [Authorization header](resources.md#authorization-header)

-> *Request body* :
>A medium's ID (pgtype.UUID)

*Example*:
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc"
}
```
-> *Error Response status code to handle* : 

    - 401 Unauthorized - Access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found - No medium with given ID found in database

-> *OK Response status code expected* :

    200 OK

-> *OK Response body example* :
>None