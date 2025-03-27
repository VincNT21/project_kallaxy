package server

import "github.com/jackc/pgx/v5/pgtype"

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
