-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, revoked_at, expires_at, user_id)
VALUES (
    $1,
    now(),
    now(),
    NULL,
    now() + interval '60 days',
    $2
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = now(), updated_at = now()
WHERE token = $1;