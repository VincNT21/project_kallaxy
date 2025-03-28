-- name: CreateMedium :one
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
RETURNING *;

-- name: UpdateMedium :one
UPDATE media
SET title = $2, creator = $3, release_year = $4, image_url = $5, metadata = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetMediumByTitle :one
SELECT * FROM media
WHERE LOWER(title) = LOWER($1);

-- name: GetMediaByType :many
SELECT * FROM media
WHERE LOWER(media_type) = LOWER($1);

-- name: DeleteMedium :one
WITH deleted AS (
    DELETE FROM media
    WHERE id = $1
)
SELECT count(*) FROM deleted;

-- name: ResetMedia :exec
DELETE FROM media;