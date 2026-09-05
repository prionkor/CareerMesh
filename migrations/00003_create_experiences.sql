-- migrations/00003_create_experiences.sql

-- +goose Up

CREATE TYPE employment_type AS ENUM (
    'full_time',
    'part_time',
    'contract',
    'freelance',
    'internship',
    'temporary',
    'volunteer'
);

CREATE TABLE experiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company TEXT NOT NULL,
    website TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL,
    employment_type employment_type NOT NULL,
    location TEXT NOT NULL DEFAULT '',
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_experiences_user_id ON experiences(user_id);

-- +goose Down

DROP TABLE experiences;

DROP TYPE employment_type;