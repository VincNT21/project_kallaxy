-- name: CreateUserMediumRecord :one
INSERT INTO users_media_records (id, created_at, updated_at, user_id, media_id, is_finished, start_date, end_date, duration)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetRecordsByUserID :many
SELECT * FROM users_media_records
WHERE user_id = $1;

-- name: UpdateRecord :one
UPDATE users_media_records
SET is_finished = $2, start_date = $3, end_date = $4, duration = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRecord :one
WITH deleted AS (
    DELETE FROM users_media_records
    WHERE id = $1
    RETURNING *
)
SELECT count(*) FROM deleted;