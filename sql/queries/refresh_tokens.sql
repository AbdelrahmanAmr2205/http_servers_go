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