-- name: StorePasswordToken :one
INSERT INTO password_reset_tokens (token, user_id, user_email, created_at, expires_at)
VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4
)
RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1;

-- name: InvalidateResetTokensByUserId :exec
UPDATE password_reset_tokens 
SET used_at = NOW()
WHERE user_id = $1;

-- name: ResetPasswordResetTable :exec
DELETE FROM password_reset_tokens;