-- +goose Up
-- +goose StatementBegin

INSERT INTO roles (code, name, description)
VALUES
    ('admin', 'Admin', 'Full access to all admin operations'),
    ('manager', 'Manager', 'Manage products, media, pages, site blocks and read stats'),
    ('content_maker', 'Content Maker', 'Manage pages and site blocks'),
    ('analyst', 'Analyst', 'Read stats and export reports')
    ON CONFLICT (code) DO NOTHING;

INSERT INTO permissions (code, description)
VALUES
    ('users.read', 'Read users'),
    ('users.write', 'Create, update, activate, deactivate users and manage roles'),

    ('products.read', 'Read products'),
    ('products.write', 'Create and update products'),
    ('products.publish', 'Publish products'),
    ('products.archive', 'Archive products'),

    ('media.read', 'Read media'),
    ('media.write', 'Attach and manage media'),
    ('media.upload', 'Upload media files'),

    ('pages.read', 'Read pages'),
    ('pages.write', 'Create and update pages'),
    ('pages.publish', 'Publish pages'),
    ('pages.archive', 'Archive pages'),

    ('site_blocks.read', 'Read site blocks'),
    ('site_blocks.write', 'Update site blocks'),

    ('stats.read', 'Read statistics'),
    ('stats.export', 'Export statistics')
    ON CONFLICT (code) DO NOTHING;


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'admin'
    ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p
              ON p.code IN (
                            'products.read',
                            'products.write',
                            'products.publish',
                            'products.archive',
                            'media.read',
                            'media.write',
                            'media.upload',
                            'pages.read',
                            'pages.write',
                            'pages.publish',
                            'pages.archive',
                            'site_blocks.read',
                            'site_blocks.write',
                            'stats.read'
                  )
WHERE r.code = 'manager'
    ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p
              ON p.code IN (
                            'pages.read',
                            'pages.write',
                            'site_blocks.read',
                            'site_blocks.write'
                  )
WHERE r.code = 'content_maker'
    ON CONFLICT DO NOTHING;

-- =========================
-- role_permissions: analyst
-- =========================
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p
              ON p.code IN (
                            'stats.read',
                            'stats.export'
                  )
WHERE r.code = 'analyst'
    ON CONFLICT DO NOTHING;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id
    FROM roles
    WHERE code IN ('admin', 'manager', 'content_maker', 'analyst')
);

DELETE FROM permissions
WHERE code IN (
               'users.read',
               'users.write',
               'products.read',
               'products.write',
               'products.publish',
               'products.archive',
               'media.read',
               'media.write',
               'media.upload',
               'pages.read',
               'pages.write',
               'pages.publish',
               'pages.archive',
               'site_blocks.read',
               'site_blocks.write',
               'stats.read',
               'stats.export'
    );

DELETE FROM roles
WHERE code IN ('admin', 'manager', 'content_maker', 'analyst');

-- +goose StatementEnd