-- migrations/00004_create_experience_bullets.sql

-- +goose Up

CREATE TABLE experience_bullets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    experience_id UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_experience_bullets_experience_id
    ON experience_bullets(experience_id);

-- +goose Down

DROP TABLE experience_bullets;