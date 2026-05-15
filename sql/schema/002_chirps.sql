-- +goose Up
-- Create a new table 'chirps' with a primary key and columns
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
-- Drop an existing table 'chirps'
DROP TABLE chirps;