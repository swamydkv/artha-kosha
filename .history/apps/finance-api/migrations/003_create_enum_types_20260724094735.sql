-- Migration: 003_create_enum_types.sql
CREATE TYPE IF NOT EXISTS session_status AS ENUM ('active','expired','revoked');
CREATE TYPE IF NOT EXISTS audit_result AS ENUM ('success','failure');
CREATE TYPE IF NOT EXISTS domain_event_type AS ENUM (
    'USER_REGISTERED', 'USER_LOGGED_IN', 'USER_LOGGED_OUT', 'ACCOUNT_CREATED', 'TRANSACTION_CREATED', 'BUDGET_CREATED'
);
CREATE TYPE IF NOT EXISTS processing_status AS ENUM ('pending','processed','failed');
