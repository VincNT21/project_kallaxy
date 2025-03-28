// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: media.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMedium = `-- name: CreateMedium :one
INSERT INTO media (id, type, created_at, updated_at, title, creator, release_year, image_url, metadata)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, type, created_at, updated_at, title, creator, release_year, image_url, metadata
`

type CreateMediumParams struct {
	Type        string
	Title       string
	Creator     string
	ReleaseYear int32
	ImageUrl    string
	Metadata    []byte
}

func (q *Queries) CreateMedium(ctx context.Context, arg CreateMediumParams) (Medium, error) {
	row := q.db.QueryRow(ctx, createMedium,
		arg.Type,
		arg.Title,
		arg.Creator,
		arg.ReleaseYear,
		arg.ImageUrl,
		arg.Metadata,
	)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Creator,
		&i.ReleaseYear,
		&i.ImageUrl,
		&i.Metadata,
	)
	return i, err
}

const deleteMedium = `-- name: DeleteMedium :one
WITH deleted AS (
    DELETE FROM media
    WHERE id = $1
)
SELECT count(*) FROM deleted
`

func (q *Queries) DeleteMedium(ctx context.Context, id pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, deleteMedium, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getMediaByType = `-- name: GetMediaByType :many
SELECT id, type, created_at, updated_at, title, creator, release_year, image_url, metadata FROM media
WHERE type = $1
`

func (q *Queries) GetMediaByType(ctx context.Context, type_ string) ([]Medium, error) {
	rows, err := q.db.Query(ctx, getMediaByType, type_)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Medium
	for rows.Next() {
		var i Medium
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Creator,
			&i.ReleaseYear,
			&i.ImageUrl,
			&i.Metadata,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMediumByTitle = `-- name: GetMediumByTitle :one
SELECT id, type, created_at, updated_at, title, creator, release_year, image_url, metadata FROM media
WHERE title = $1
`

func (q *Queries) GetMediumByTitle(ctx context.Context, title string) (Medium, error) {
	row := q.db.QueryRow(ctx, getMediumByTitle, title)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Creator,
		&i.ReleaseYear,
		&i.ImageUrl,
		&i.Metadata,
	)
	return i, err
}

const resetMedia = `-- name: ResetMedia :exec
DELETE FROM media
`

func (q *Queries) ResetMedia(ctx context.Context) error {
	_, err := q.db.Exec(ctx, resetMedia)
	return err
}

const updateMedium = `-- name: UpdateMedium :one
UPDATE media
SET title = $2, creator = $3, release_year = $4, image_url = $5, metadata = $6, updated_at = NOW()
WHERE id = $1
RETURNING id, type, created_at, updated_at, title, creator, release_year, image_url, metadata
`

type UpdateMediumParams struct {
	ID          pgtype.UUID
	Title       string
	Creator     string
	ReleaseYear int32
	ImageUrl    string
	Metadata    []byte
}

func (q *Queries) UpdateMedium(ctx context.Context, arg UpdateMediumParams) (Medium, error) {
	row := q.db.QueryRow(ctx, updateMedium,
		arg.ID,
		arg.Title,
		arg.Creator,
		arg.ReleaseYear,
		arg.ImageUrl,
		arg.Metadata,
	)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Creator,
		&i.ReleaseYear,
		&i.ImageUrl,
		&i.Metadata,
	)
	return i, err
}
