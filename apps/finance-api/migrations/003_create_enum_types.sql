-- 003_create_enum_types.sql
-- Create enum types used by sessions, audit, domain events, and processing status

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'session_status') THEN
        CREATE TYPE session_status AS ENUM ('active','expired','revoked');
    END IF;
EXCEPTION WHEN others THEN
    -- ignore if already exists
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'audit_result') THEN
        CREATE TYPE audit_result AS ENUM ('success','failure');
    END IF;
EXCEPTION WHEN others THEN
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'domain_event_type') THEN
        CREATE TYPE domain_event_type AS ENUM ('USER_REGISTERED','USER_LOGGED_IN','USER_LOGGED_OUT','ACCOUNT_CREATED','TRANSACTION_CREATED','BUDGET_CREATED');
    END IF;
EXCEPTION WHEN others THEN
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing_status') THEN
        CREATE TYPE processing_status AS ENUM ('pending','processed','failed');
    END IF;
EXCEPTION WHEN others THEN
END$$;

