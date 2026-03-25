-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         citext      NOT NULL,
    password_hash text        NOT NULL,
    is_active     boolean     NOT NULL DEFAULT true,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT users_email_uk UNIQUE (email)
    );

CREATE INDEX IF NOT EXISTS users_is_active_idx
    ON users(is_active);

DROP TRIGGER IF EXISTS users_set_updated_at ON users;

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS roles (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        varchar(64)  NOT NULL,
    name        varchar(128) NOT NULL,
    description text         NULL,

    CONSTRAINT roles_code_uk UNIQUE (code)
    );

CREATE TABLE IF NOT EXISTS permissions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        varchar(128) NOT NULL,
    description text         NULL,

    CONSTRAINT permissions_code_uk UNIQUE (code)
    );

CREATE TABLE IF NOT EXISTS user_roles (
    user_id    uuid        NOT NULL,
    role_id    uuid        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT user_roles_pk PRIMARY KEY (user_id, role_id),

    CONSTRAINT user_roles_user_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT user_roles_role_fk
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS user_roles_user_idx
    ON user_roles(user_id);

CREATE INDEX IF NOT EXISTS user_roles_role_idx
    ON user_roles(role_id);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       uuid NOT NULL,
    permission_id uuid NOT NULL,

    CONSTRAINT role_permissions_pk PRIMARY KEY (role_id, permission_id),

    CONSTRAINT role_permissions_role_fk
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,

    CONSTRAINT role_permissions_permission_fk
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS role_permissions_role_idx
    ON role_permissions(role_id);

CREATE INDEX IF NOT EXISTS role_permissions_permission_idx
    ON role_permissions(permission_id);

CREATE TABLE IF NOT EXISTS refresh_sessions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL,
    token_hash  text        NOT NULL,
    expires_at  timestamptz NOT NULL,
    revoked_at  timestamptz NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT refresh_sessions_user_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT refresh_sessions_token_hash_uk UNIQUE (token_hash),

    CONSTRAINT refresh_sessions_expires_check
    CHECK (expires_at > created_at)
    );

CREATE INDEX IF NOT EXISTS refresh_sessions_user_idx
    ON refresh_sessions(user_id);

CREATE INDEX IF NOT EXISTS refresh_sessions_expires_idx
    ON refresh_sessions(expires_at);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL,
    token_hash  text        NOT NULL,
    expires_at  timestamptz NOT NULL,
    used_at     timestamptz NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT password_reset_tokens_user_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT password_reset_tokens_token_hash_uk UNIQUE (token_hash),

    CONSTRAINT password_reset_tokens_expires_check
    CHECK (expires_at > created_at)
    );

CREATE INDEX IF NOT EXISTS password_reset_tokens_user_idx
    ON password_reset_tokens(user_id);

CREATE INDEX IF NOT EXISTS password_reset_tokens_expires_idx
    ON password_reset_tokens(expires_at);

CREATE UNIQUE INDEX IF NOT EXISTS password_reset_tokens_one_active_per_user
    ON password_reset_tokens(user_id)
    WHERE used_at IS NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS password_reset_tokens_one_active_per_user;

DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS refresh_sessions;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS set_updated_at();


-- +goose StatementEnd