CREATE TABLE users (
    id UUID PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    roles JSONB NOT NULL,
    last_login_at TIMESTAMPTZ
);