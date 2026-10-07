-- migrations/00015_seed_roles.sql

-- +goose Up

INSERT INTO roles (name)
VALUES ('user'), ('admin'), ('superadmin')
ON CONFLICT (name) DO NOTHING;

-- +goose Down

DELETE FROM roles
WHERE name IN ('user', 'admin', 'superadmin');
