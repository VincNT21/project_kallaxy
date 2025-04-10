-- name: CreateUserMediumRecord :one
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
RETURNING *;

-- name: GetRecordsByUserID :many
SELECT * FROM users_media_records
WHERE user_id = $1;

-- name: GetRecordsAndMediaByUserID :many
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
WHERE records.user_id = $1;

-- name: GetRecordByID :one
SELECT * FROM users_media_records
WHERE id = $1;

-- name: UpdateRecord :one
UPDATE users_media_records
SET is_finished = $2, start_date = $3, end_date = $4, duration = $5, comments = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRecord :one
WITH deleted AS (
    DELETE FROM users_media_records
    WHERE id = $1
    RETURNING *
)
SELECT count(*) FROM deleted;

-- name: GetDatesFromRecord :one
SELECT start_date, end_date
FROM users_media_records
WHERE id = $1;

-- name: ResetRecords :exec
DELETE FROM users_media_records;