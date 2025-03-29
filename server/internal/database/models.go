// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Medium struct {
	ID          pgtype.UUID
	MediaType   string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	Title       string
	Creator     string
	ReleaseYear int32
	ImageUrl    pgtype.Text
	Metadata    []byte
}

type PasswordResetToken struct {
	Token     string
	UserID    pgtype.UUID
	UserEmail string
	CreatedAt pgtype.Timestamp
	ExpiresAt pgtype.Timestamp
	UsedAt    pgtype.Timestamp
}

type RefreshToken struct {
	Token     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	UserID    pgtype.UUID
	ExpiresAt pgtype.Timestamp
	RevokedAt pgtype.Timestamp
}

type User struct {
	ID             pgtype.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
	Username       string
	HashedPassword string
	Email          string
}

type UsersMediaRecord struct {
	ID         pgtype.UUID
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
	UserID     pgtype.UUID
	MediaID    pgtype.UUID
	IsFinished pgtype.Bool
	StartDate  pgtype.Timestamp
	EndDate    pgtype.Timestamp
	Duration   pgtype.Interval
}
