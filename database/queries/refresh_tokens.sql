-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :one
WITH revoked AS (
    UPDATE refresh_tokens
    SET revoked_at = NOW(), updated_at = NOW()
    WHERE token = $1
    RETURNING *
)
SELECT count(*) FROM revoked;

-- name: RevokeAllRefreshTokensByUserID :one
WITH revoked AS (
    UPDATE refresh_tokens
    SET revoked_at = NOW(), updated_at = NOW()
    WHERE user_id = $1
    RETURNING *
)
SELECT count(*) FROM revoked;

-- name: DeleteRevokedTokens :exec
DELETE FROM refresh_tokens
WHERE revoked_at IS NOT NULL;