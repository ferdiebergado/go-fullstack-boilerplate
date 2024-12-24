-- Active: 1735025158779@@127.0.0.1@5432@gfb
CREATE TYPE auth_method AS ENUM ('email/password', 'oauth');

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid () PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    oauth_provider VARCHAR(50),
    oauth_id VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    auth_method auth_method NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);