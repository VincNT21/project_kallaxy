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
	MediaType   string           `json:"media_type"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	Title       string           `json:"title"`
	Creator     string           `json:"creator"`
	ReleaseYear int32            `json:"release_year"`
	ImageUrl    pgtype.Text      `json:"image_url"`
	Metadata    json.RawMessage  `json:"metadata"`
}

type Record struct {
	ID         pgtype.UUID      `json:"id"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
	UserID     pgtype.UUID      `json:"user_id"`
	MediaID    pgtype.UUID      `json:"media_id"`
	IsFinished pgtype.Bool      `json:"is_finished"`
	StartDate  pgtype.Timestamp `json:"start_date"`
	EndDate    pgtype.Timestamp `json:"end_date"`
	Duration   int32            `json:"duration"`
}

type MediumWithRecord struct {
	ID          pgtype.UUID      `json:"record_id"`
	UserID      pgtype.UUID      `json:"user_id"`
	MediaID     pgtype.UUID      `json:"medium_id"`
	IsFinished  pgtype.Bool      `json:"is_finished"`
	StartDate   pgtype.Timestamp `json:"start_date"`
	EndDate     pgtype.Timestamp `json:"end_date"`
	Duration    int32            `json:"duration"`
	MediaType   string           `json:"media_type"`
	Title       string           `json:"title"`
	Creator     string           `json:"creator"`
	ReleaseYear int32            `json:"release_year"`
	ImageUrl    pgtype.Text      `json:"image_url"`
	Metadata    []byte           `json:"metadata"`
}
