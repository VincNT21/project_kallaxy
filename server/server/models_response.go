package server

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Medium struct {
	ID          pgtype.UUID      `json:"id"`
	Type        string           `json:"type"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	Title       string           `json:"title"`
	Creator     string           `json:"creator"`
	ReleaseYear int32            `json:"release_year"`
	ImageUrl    string           `json:"image_url"`
	Metadata    json.RawMessage  `json:"metadata"`
}
