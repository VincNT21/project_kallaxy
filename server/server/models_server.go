package server

import (
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
	ID          pgtype.UUID            `json:"id"`
	MediaType   string                 `json:"media_type"`
	CreatedAt   pgtype.Timestamp       `json:"created_at"`
	UpdatedAt   pgtype.Timestamp       `json:"updated_at"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	ReleaseYear string                 `json:"release_year"`
	ImageUrl    string                 `json:"image_url"`
	Metadata    map[string]interface{} `json:"metadata"`
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
	Comments   string           `json:"comments"`
}

type MediumWithRecord struct {
	ID          pgtype.UUID            `json:"record_id"`
	UserID      pgtype.UUID            `json:"user_id"`
	MediaID     pgtype.UUID            `json:"medium_id"`
	IsFinished  pgtype.Bool            `json:"is_finished"`
	StartDate   pgtype.Timestamp       `json:"start_date"`
	EndDate     pgtype.Timestamp       `json:"end_date"`
	Duration    int32                  `json:"duration"`
	Comments    string                 `json:"comments"`
	MediaType   string                 `json:"media_type"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	ReleaseYear string                 `json:"release_year"`
	ImageUrl    string                 `json:"image_url"`
	Metadata    map[string]interface{} `json:"metadata"`
}
