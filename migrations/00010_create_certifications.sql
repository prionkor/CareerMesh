-- migrations/00010_create_certifications.sql

-- +goose Up

CREATE TABLE certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    issuer TEXT NOT NULL DEFAULT '',
    issue_date DATE,
    expiry_date DATE,
    credential_id TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_certifications_user_id ON certifications(user_id);

-- +goose Down

DROP TABLE certifications;