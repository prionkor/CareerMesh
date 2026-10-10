-- +goose Up

CREATE TABLE job_descriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    company_website TEXT NOT NULL DEFAULT '',
    source_url TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    employment_type employment_type,
    original_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_descriptions_user_id ON job_descriptions(user_id);

-- +goose Down

DROP TABLE job_descriptions;