CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    roles JSONB NOT NULL,
    last_login_at TIMESTAMPTZ
);