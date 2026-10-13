-- migrations/00016_create_job_descriptions.sql
-- +goose Up

CREATE TABLE job_descriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_version_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_descriptions_user_id
    ON job_descriptions(user_id);

-- +goose Down

DROP TABLE job_descriptions;