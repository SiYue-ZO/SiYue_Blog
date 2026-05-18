-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ENUM types
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');
CREATE TYPE post_status AS ENUM ('draft', 'published', 'archived');

-- Auto-update updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- users table
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)   NOT NULL,
    email         VARCHAR(255)  NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,
    avatar_url    VARCHAR(512),
    role          VARCHAR(20)   NOT NULL DEFAULT 'user',
    status        user_status   NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_username ON users (username);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
