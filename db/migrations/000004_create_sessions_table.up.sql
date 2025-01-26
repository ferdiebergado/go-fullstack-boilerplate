CREATE TABLE sessions (
    id BYTEA PRIMARY KEY, -- session ID
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Foreign key with ON DELETE CASCADE
    ip_address INET NOT NULL,                   -- PostgreSQL's INET type for IP addresses
    user_agent TEXT NOT NULL,
    login_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Timestamp with timezone
    last_activity TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expiry_time TIMESTAMPTZ,           -- Optional expiry time
    data JSONB                 -- PostgreSQL's JSONB for efficient JSON storage and querying
);

CREATE INDEX idx_user_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_user_sessions_last_activity ON sessions (last_activity);
CREATE INDEX idx_user_sessions_expiry_time ON sessions (expiry_time) WHERE expiry_time IS NOT NULL; -- Partial index
