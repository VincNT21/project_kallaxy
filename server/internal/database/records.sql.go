// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: records.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUserMediumRecord = `-- name: CreateUserMediumRecord :one
INSERT INTO users_media_records (id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments
`

type CreateUserMediumRecordParams struct {
	UserID     pgtype.UUID
	MediaID    pgtype.UUID
	IsFinished pgtype.Bool
	StartDate  pgtype.Timestamp
	EndDate    pgtype.Timestamp
	Duration   pgtype.Interval
	Comments   string
}

func (q *Queries) CreateUserMediumRecord(ctx context.Context, arg CreateUserMediumRecordParams) (UsersMediaRecord, error) {
	row := q.db.QueryRow(ctx, createUserMediumRecord,
		arg.UserID,
		arg.MediaID,
		arg.IsFinished,
		arg.StartDate,
		arg.EndDate,
		arg.Duration,
		arg.Comments,
	)
	var i UsersMediaRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.MediaID,
		&i.IsFinished,
		&i.StartDate,
		&i.EndDate,
		&i.Duration,
		&i.Comments,
	)
	return i, err
}

const deleteRecord = `-- name: DeleteRecord :one
WITH deleted AS (
    DELETE FROM users_media_records
    WHERE media_id = $1
    AND user_id = $2
    RETURNING id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments
)
SELECT count(*) FROM deleted
`

type DeleteRecordParams struct {
	MediaID pgtype.UUID
	UserID  pgtype.UUID
}

func (q *Queries) DeleteRecord(ctx context.Context, arg DeleteRecordParams) (int64, error) {
	row := q.db.QueryRow(ctx, deleteRecord, arg.MediaID, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDatesFromRecord = `-- name: GetDatesFromRecord :one
SELECT start_date, end_date
FROM users_media_records
WHERE id = $1
`

type GetDatesFromRecordRow struct {
	StartDate pgtype.Timestamp
	EndDate   pgtype.Timestamp
}

func (q *Queries) GetDatesFromRecord(ctx context.Context, id pgtype.UUID) (GetDatesFromRecordRow, error) {
	row := q.db.QueryRow(ctx, getDatesFromRecord, id)
	var i GetDatesFromRecordRow
	err := row.Scan(&i.StartDate, &i.EndDate)
	return i, err
}

const getRecordByID = `-- name: GetRecordByID :one
SELECT id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments FROM users_media_records
WHERE id = $1
`

func (q *Queries) GetRecordByID(ctx context.Context, id pgtype.UUID) (UsersMediaRecord, error) {
	row := q.db.QueryRow(ctx, getRecordByID, id)
	var i UsersMediaRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.MediaID,
		&i.IsFinished,
		&i.StartDate,
		&i.EndDate,
		&i.Duration,
		&i.Comments,
	)
	return i, err
}

const getRecordsAndMediaByUserID = `-- name: GetRecordsAndMediaByUserID :many
SELECT
    records.id, 
    records.user_id, 
    records.media_id, 
    records.is_finished, 
    records.start_date, 
    records.end_date, 
    records.duration, 
    records.comments,
    media.media_type,
    media.title,
    media.creator,
    media.pub_date,
    media.image_url,
    media.metadata
FROM users_media_records AS records
INNER JOIN media
ON records.media_id = media.id
WHERE records.user_id = $1
`

type GetRecordsAndMediaByUserIDRow struct {
	ID         pgtype.UUID
	UserID     pgtype.UUID
	MediaID    pgtype.UUID
	IsFinished pgtype.Bool
	StartDate  pgtype.Timestamp
	EndDate    pgtype.Timestamp
	Duration   pgtype.Interval
	Comments   string
	MediaType  string
	Title      string
	Creator    string
	PubDate    string
	ImageUrl   string
	Metadata   []byte
}

func (q *Queries) GetRecordsAndMediaByUserID(ctx context.Context, userID pgtype.UUID) ([]GetRecordsAndMediaByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getRecordsAndMediaByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRecordsAndMediaByUserIDRow
	for rows.Next() {
		var i GetRecordsAndMediaByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.MediaID,
			&i.IsFinished,
			&i.StartDate,
			&i.EndDate,
			&i.Duration,
			&i.Comments,
			&i.MediaType,
			&i.Title,
			&i.Creator,
			&i.PubDate,
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

const getRecordsByUserID = `-- name: GetRecordsByUserID :many
SELECT id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments FROM users_media_records
WHERE user_id = $1
`

func (q *Queries) GetRecordsByUserID(ctx context.Context, userID pgtype.UUID) ([]UsersMediaRecord, error) {
	rows, err := q.db.Query(ctx, getRecordsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersMediaRecord
	for rows.Next() {
		var i UsersMediaRecord
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.MediaID,
			&i.IsFinished,
			&i.StartDate,
			&i.EndDate,
			&i.Duration,
			&i.Comments,
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

const resetRecords = `-- name: ResetRecords :exec
DELETE FROM users_media_records
`

func (q *Queries) ResetRecords(ctx context.Context) error {
	_, err := q.db.Exec(ctx, resetRecords)
	return err
}

const updateRecord = `-- name: UpdateRecord :one
UPDATE users_media_records
SET is_finished = $2, start_date = $3, end_date = $4, duration = $5, comments = $6, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration, comments
`

type UpdateRecordParams struct {
	ID         pgtype.UUID
	IsFinished pgtype.Bool
	StartDate  pgtype.Timestamp
	EndDate    pgtype.Timestamp
	Duration   pgtype.Interval
	Comments   string
}

func (q *Queries) UpdateRecord(ctx context.Context, arg UpdateRecordParams) (UsersMediaRecord, error) {
	row := q.db.QueryRow(ctx, updateRecord,
		arg.ID,
		arg.IsFinished,
		arg.StartDate,
		arg.EndDate,
		arg.Duration,
		arg.Comments,
	)
	var i UsersMediaRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.MediaID,
		&i.IsFinished,
		&i.StartDate,
		&i.EndDate,
		&i.Duration,
		&i.Comments,
	)
	return i, err
}
