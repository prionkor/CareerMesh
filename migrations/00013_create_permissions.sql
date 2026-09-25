-- migrations/00013_create_permissions.sql

-- +goose Up

-- Permission keys are stable strings rather than an enum so new permissions
-- can be added without a schema migration.
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE permissions;
