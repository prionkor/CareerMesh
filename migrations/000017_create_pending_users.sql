-- migrations/000017_create_pending_users.sql

-- +goose Up
CREATE TABLE pending_users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pending_users_token_hash ON pending_users(token_hash);

-- +goose Down
DROP TABLE pending_users;
