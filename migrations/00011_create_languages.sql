-- migrations/00011_create_languages.sql

-- +goose Up

CREATE TYPE language_proficiency AS ENUM (
    'beginner',
    'elementary',
    'intermediate',
    'upper_intermediate',
    'advanced',
    'native'
);

CREATE TABLE languages (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    proficiency language_proficiency NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_languages_user_id
    ON languages(user_id);

-- +goose Down

DROP TABLE languages;

DROP TYPE language_proficiency;