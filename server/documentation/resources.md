# Project Kallaxy Resources <!-- omit from toc -->

## Headers

### Authorization header
Most endpoint needs a valid access token, some needs a valid refresh token.
This token must be set in an "Authorization" header.

Access token is in format 
Refresh token is in format 

```json
{
    "Authorization": "Bearer <token>"
}
```

## User resource

### Structure
- `id`:         *string* (UUIDv4 format) - User's unique identifier
- `created_at`: *string* (ISO 8601 datetime) - When the user was created
- `updated_at`: *string* (ISO 8601 datetime) - Last time the user's info was updated
- `username`:   *string* - User's chosen username
- `email`:      *string* - User's email adress
  
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

### GO
```go
type User struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
}
```

## Media resource

### Structure
- `id`:             *string* (UUIDv4 format) - Medium's unique identifier
- `media_type`:     *string* - Medium's type (book, movie, serie...)
- `created_at`:     *string* (ISO 8601 datetime format) - When the medium was first created
- `updated_at`:     *string* (ISO 8601 datetime format) - Last time the medium's info was updated
- `title`:          *string* - Medium's title
- `creator`:        *string* - Medium's creator (author, director...)
- `release_year`:   *int32* - Medium's year of publication
- `image_url`:      *string* - A link to medium's cover
- `metadata`:       *json.RawMessage* - A json object containing metatadata about the medium, according to media type (see below)

### Example
```json
{
    "id": "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc",
    "type": "book",
    "created_at": "2025-03-26T14:20:23.525332",
    "updated_at": "2025-03-26T14:20:23.525332",
    "title": "The Fellowship of the ring",
    "creator": "J.R.R. Tolkien",
    "release_year": 1954,
    "image_url": "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
    "metadata": ""
}
```

## Metadata for media


## Record resource 

### Structure
- `id`:             *string* (UUIDv4 format) - Record's unique identifier
- `created_at`:     *string* (ISO 8601 datetime) - When the record was first created
- `updated_at`:     *string* (ISO 8601 datetime) - Last time the user info was updated
- `user_id`:        *string* (UUIDv4 format) - User concerned by the record
- `media_id`:       *string* (UUIDv4 format) - Medium concerned by the record
- `is_finished`:    *boolean* - Does user have finished reading/watching/playing the medium
- `start_date`:     *string* (ISO 8601 datetime) - When user started to read/watch/play the medium
- `end_date`:       *string* (ISO 8601 datetime) - When user finished reading/watching/playing the medium
- `duration`:       *int32* - Auto-calculated days interval between start and end dates

### Example
```json
{
    "id": "4aea83e5-36e2-47c3-a121-7e3db9ac72d1",
    "created_at": "2025-03-31T08:59:09.523473",
    "updated_at": "2025-03-31T08:59:09.523473",
    "user_id": "2a0d54f8-37b8-4e51-826d-6f9632c374a4",
    "media_id": "3b75af06-e596-42ce-a953-bf235dfc9102",
    "is_finished": true,
    "start_date": "2025-03-26T14:20:23.525332",
    "end_date": "2025-03-31T08:47:29.205805",
    "duration": 4
}
```

## Specific formats
### Tokens
#### Access token
A JSON Web Token consisting of three parts separated by dots (header.payload.signature)
Example: 
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrYWxsYXh5Iiwic3ViIjoiODFjMWNiMGQtYmJkYi00ZmFhLWFlZGUtYmQzNzFhNGFiNzIyIiwiZXhwIjoxNzQzMDIxMTMyLCJpYXQiOjE3NDMwMTc1MzJ9.1PUE_93e6pXaLwjZiMIfr5DAcxTxE4jEIiRftQuJptI"
}

```

Access token expiration time : 1 hour

#### Refresh token
A 64-character hexadecimal string used to obtain a new access token
Example: "176ddabd5f4c932b8cda583e00b620a05242187680002a071e8a13c4e2b0b14f"

Refresh token lifspan : 30 days

### UUID
A UUID string in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" where x is a hexadecimal digit.
Example: "d8b5ad72-1a8d-4990-bb83-44bd4daa32dc"

### Datetime
A timestamp string in RFC3339 format: YYYY-MM-DDThh:mm:ssZ
Where:
- YYYY-MM-DD is the date portion
- T is a literal character separating date and time
- hh:mm:ss.sss is the time with optional microsecond precision
- Z indicates UTC timezone (can be replaced with +/-hh:mm offset)

Examples:
  - "2025-04-01T07:58:56Z" (basic format)
  - "2025-04-01T07:58:56.827795Z" (with microsecond precision)

All timestamps should be in UTC timezone (indicated by the 'Z' suffix).

In Go, you can use 
```go
time.Now().UTC().Format(time.RFC3339)
```

## Go models

### For client's requests

#### Users
```go
type parametersCreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
```

```go
type parametersUpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
```

#### Media
```go
type parametersCreateMedium struct {
	Title       string          `json:"title"`
	MediaType   string          `json:"media_type"`
	Creator     string          `json:"creator"`
	ReleaseYear int32           `json:"release_year"`
	ImageUrl    string          `json:"image_url"`
	Metadata    json.RawMessage `json:"metadata"`
}
```

```go
type parametersUpdateMedium struct {
	MediumID    string          `json:"medium_id"`
	Title       string          `json:"title"`
	Creator     string          `json:"creator"`
	ReleaseYear int32           `json:"release_year"`
	ImageUrl    string          `json:"image_url"`
	Metadata    json.RawMessage `json:"metadata"`
}
```

```go
type parametersDeleteMedium struct {
	MediumID string `json:"medium_id"`
}
```

#### Records
```go
type parametersCreateUserMediumRecord struct {
	MediumID  string `json:"medium_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
```

```go
type parametersUpdateRecord struct {
	RecordID  string `json:"record_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
```

```go
type parametersDeleteRecord struct {
	RecordID string `json:"record_id"`
}
```

#### Authentification

```go
type parametersLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```

#### Admin/Password Reset
```go
type parametersPasswordResetRequest struct {
	Email string `json:"email"`
}
```


### For server's responses
```go
type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}
```

```go
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
```

```go
type TokensAndUser struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
```

```go
type Medium struct {
	ID          string          `json:"id"`
	MediaType   string          `json:"media_type"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Title       string          `json:"title"`
	Creator     string          `json:"creator"`
	ReleaseYear int32           `json:"release_year"`
	ImageUrl    string          `json:"image_url"`
	Metadata    json.RawMessage `json:"metadata"`
}
```

```go
type ListMedia struct {
	Media []Medium `json:"media"`
}
```

```go
type Record struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	UserID     string `json:"user_id"`
	MediaID    string `json:"media_id"`
	IsFinished bool   `json:"is_finished"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Duration   int32  `json:"duration"`
}
```

```go
type Records struct {
	Records []Record `json:"records"`
}
```

#### Admin-Password Reset
```go
type PasswordResetRequest struct {
	Message    string `json:"message"`
	ResetLink  string `json:"reset_link"`
	ResetToken string `json:"reset_token"`
}
```

```go
type responseVerifyResetToken struct {
	Valid bool   `json:"valid"`
	Email string `json:"email"`
}
```