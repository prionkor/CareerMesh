-- migrations/00009_create_education.sql

-- +goose Up

CREATE TABLE education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL DEFAULT '',
    field TEXT NOT NULL DEFAULT '',
    start_date DATE,
    end_date DATE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_education_user_id ON education(user_id);

-- +goose Down

DROP TABLE education;