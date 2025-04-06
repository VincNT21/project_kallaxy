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
INSERT INTO media (id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata)
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
RETURNING id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata
`

type CreateMediumParams struct {
	MediaType   string
	Title       string
	Creator     string
	ReleaseYear int32
	ImageUrl    pgtype.Text
	Metadata    []byte
}

func (q *Queries) CreateMedium(ctx context.Context, arg CreateMediumParams) (Medium, error) {
	row := q.db.QueryRow(ctx, createMedium,
		arg.MediaType,
		arg.Title,
		arg.Creator,
		arg.ReleaseYear,
		arg.ImageUrl,
		arg.Metadata,
	)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaType,
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
    RETURNING id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata
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
SELECT id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata FROM media
WHERE LOWER(media_type) = LOWER($1)
`

func (q *Queries) GetMediaByType(ctx context.Context, lower string) ([]Medium, error) {
	rows, err := q.db.Query(ctx, getMediaByType, lower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Medium
	for rows.Next() {
		var i Medium
		if err := rows.Scan(
			&i.ID,
			&i.MediaType,
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

const getMediumByID = `-- name: GetMediumByID :one
SELECT id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata FROM media
WHERE id = $1
`

func (q *Queries) GetMediumByID(ctx context.Context, id pgtype.UUID) (Medium, error) {
	row := q.db.QueryRow(ctx, getMediumByID, id)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaType,
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

const getMediumByTitleAndType = `-- name: GetMediumByTitleAndType :one
SELECT id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata FROM media
WHERE LOWER(title) = LOWER($1)
AND LOWER(media_type) = LOWER($2)
`

type GetMediumByTitleAndTypeParams struct {
	Lower   string
	Lower_2 string
}

func (q *Queries) GetMediumByTitleAndType(ctx context.Context, arg GetMediumByTitleAndTypeParams) (Medium, error) {
	row := q.db.QueryRow(ctx, getMediumByTitleAndType, arg.Lower, arg.Lower_2)
	var i Medium
	err := row.Scan(
		&i.ID,
		&i.MediaType,
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
RETURNING id, media_type, created_at, updated_at, title, creator, release_year, image_url, metadata
`

type UpdateMediumParams struct {
	ID          pgtype.UUID
	Title       string
	Creator     string
	ReleaseYear int32
	ImageUrl    pgtype.Text
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
		&i.MediaType,
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
