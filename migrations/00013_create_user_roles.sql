-- migrations/00013_create_user_roles.sql

-- +goose Up

CREATE TABLE user_roles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_roles_role_id
    ON user_roles(role_id);

-- +goose Down

DROP TABLE user_roles;