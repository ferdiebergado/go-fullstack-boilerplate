DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_user_sessions_last_activity;
DROP INDEX IF EXISTS idx_user_sessions_expiry_time;

DROP TABLE IF EXISTS user_sessions;
