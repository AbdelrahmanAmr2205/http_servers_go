-- +goose Up
CREATE TABLE refresh_tokens (
    token CHAR(32) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    expires_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP -- Drop an existing table 'refresh_tokens'
DROP TABLE refresh_tokens;