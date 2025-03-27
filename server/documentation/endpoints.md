# Project Kallaxy Endpoints <!-- omit from toc -->

## Public API Endpoints

### Users related Endpoints

#### POST /api/users -- User creation
-> *Description* : 

    Create a new user in **users** table

-> *Request need* : 
>A username, a password and an email (string)

-> *Request body example* :
```json
{
    "username": "VincNT21",
    "password": "12345abcde",
    "email": "vincnt21@example.com"   
}

```
-> *OK Response status code expected* : 

    201-Created

-> *Error Response status code to handle* : 

    - 409-Conflict : username or email is already used by another user 

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