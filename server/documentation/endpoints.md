# Project Kallaxy Endpoints <!-- omit from toc -->


## 1. Users related Endpoints

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
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "username": "VincNT21",
    "email": "vincnt21@example.com"
}
```

### 1.2. PUT /api/users -- User info update
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
>A username (string), a password (string) and a email (string). If a field is not updated, client still need to send old info (no comparison is done in server, the three fields are replaced).

*Example*:
```json
{
    "username": "VincNT21",
    "password": "12345ghjk",
    "email": "vincnt21@example.com" 
}

```
-> *Error Response status code to handle* : 

    - 400 Bad Request : one or many fields are missing in request
    - 401 Unauthorized : means that access token is expired, client should fetch **POST /auth/refresh** to get a new access token

-> *OK Response status code expected* : 

    200 OK
    /!\ If client receives 200 OK response, it also means that user has been logged out and need to log in again for new refresh token /!\ 

-> *Response body example* :
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "username": "VincNT21",
    "email": "vincnt21@example.com"
}
```


### 1.3. POST /auth/login -- Authentification
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

    - 401 Unauthorized: given username/password does not match.

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

### 1.4. POST /auth/logout -- Logout a user
-> *Description* : 
> Logout a user by revoking all their refresh tokens

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

    - 401 Unauthorized : means that access token is expired, client should fetch **POST /auth/refresh** to get a new access token
    - 404 Not Found: No refresh token associated to user's ID (from access token) found in database

-> *OK Response status code expected* : 

    204 No Content

### 1.5. POST /auth/refresh -- Refresh access token
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

    - 401 Unauthorized: Refresh token doesn't exist in server's database or has been revoked or has expired. Client should fetch **POST /auth/login** to get a new refresh token.

-> *OK Response status code expected* : 

    201 Created

-> *Response body example* :
```json
{
    "access_token": "<access_token>",
    "refresh_token": "<refresh_token>"
}
```

### 1.6. POST /auth/revoke -- Revoke a refresh token
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

    - 401 Unauthorized: There is a problem with "Authorization" header: missing or malformed
    - 404 Not Found: The refresh token doesn't exist in server's database

-> *OK Response status code expected* : 

    204 No Content



-> *Response body example* :
>Empty

## 2. Model
-> *Description* : 
>

-> *Request headers* : 
>

*Example*:
```json
{
      
}

``` 

-> *Request body* :
>

*Example*:
```json
{
      
}

```
-> *OK Response status code expected* : 

    200

-> *Error Response status code to handle* : 

    - 400

-> *Response body example* :
```json
{
    
}
```