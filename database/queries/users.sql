-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, hashed_password, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET username = $2, hashed_password = $3, email = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET hashed_password = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
WITH deleted AS (
    DELETE FROM users
    WHERE id = $1
)
SELECT count(*) FROM deleted;

-- name: ResetUsers :exec
DELETE FROM users;