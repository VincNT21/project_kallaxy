# Project Kallaxy Resources <!-- omit from toc -->

## Headers

### Authorization header
Most endpoint needs a valid access token, some needs a valid refresh token.
This token must be set in an "Authorization" header.

```json
{
    "Authorization": "Bearer <token>"
}
```

## User resource

### Structure
- `id`: string - User's unique identifier
- `created_at`: string (ISO 8601 datetime) - When the user was created
- `updated_at`: string (ISO 8601 datetime) - Last time the user's info was updated
- `username`: string - User's chosen username
- `email`: string - User's email adress
  
### Example
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "username": "VincNT21",
    "email": "vincnt21@example.com"
}
```

## Media resource

### Structure
- `id`: string - Medium's unique identifier
- `type`: string - Medium's type (book, movie, serie...)
- `created_at`: string (ISO 8601 datetime) - When the user was created
- `updated_at`: string (ISO 8601 datetime) - Last time the user info was updated
- `title`: string - Medium's title
- `creator`: string - Medium's creator (author, director...)
- `release_year`: int32 - Medium's year of publication
- `image_url`: string - a link to medium's cover
- `metadata`: json.RawMessage - a json object, according to media type (see below)

### Example
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "type": "book",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "title": "The Fellowship of the ring",
    "creator": "J.R.R. Tolkien",
    "release_year": "1954",
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```