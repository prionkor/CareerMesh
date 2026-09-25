-- migrations/00016_seed_rbac.sql

-- +goose Up

INSERT INTO permissions (key) VALUES
    ('users:read'),
    ('users:update'),
    ('users:delete'),
    ('profiles:read'),
    ('profiles:create'),
    ('profiles:update'),
    ('profiles:delete'),
    ('experiences:read'),
    ('experiences:create'),
    ('experiences:update'),
    ('experiences:delete'),
    ('projects:read'),
    ('projects:create'),
    ('projects:update'),
    ('projects:delete'),
    ('skills:read'),
    ('skills:create'),
    ('skills:update'),
    ('skills:delete'),
    ('education:read'),
    ('education:create'),
    ('education:update'),
    ('education:delete'),
    ('certifications:read'),
    ('certifications:create'),
    ('certifications:update'),
    ('certifications:delete'),
    ('languages:read'),
    ('languages:create'),
    ('languages:update'),
    ('languages:delete'),
    -- Grants ownership-bypass for resource operations; see internal/authorization.
    ('system:admin');

INSERT INTO roles (name) VALUES
    ('user'),
    ('admin');

-- The 'user' role gets every resource permission. Ownership is still
-- enforced separately in the service layer for individual resources.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'user' AND p.key <> 'system:admin';

-- The 'admin' role gets every permission, including the ownership bypass.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin';

-- Backfill the 'user' role onto any users created before RBAC existed.
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE r.name = 'user'
ON CONFLICT DO NOTHING;

-- +goose Down

DELETE FROM user_roles;
DELETE FROM role_permissions;
DELETE FROM roles WHERE name IN ('user', 'admin');
DELETE FROM permissions;
