-- 008_gdpr_archiving.sql

CREATE TABLE IF NOT EXISTS archived_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_user_id UUID NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    mobile VARCHAR(20),
    dob DATE,
    password_hash TEXT NOT NULL,
    archived_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    retention_expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_archived_users_original_id ON archived_users(original_user_id);
CREATE INDEX IF NOT EXISTS idx_archived_users_retention ON archived_users(retention_expires_at);

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS is_archived BOOLEAN DEFAULT FALSE;

ALTER TYPE domain_event_type ADD VALUE IF NOT EXISTS 'USER_DELETED';