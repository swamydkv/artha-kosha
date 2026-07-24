-- Migration: Create users table
-- Description: Core user account storage for authentication

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    mobile_number TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_username CHECK (username ~ '^[A-Za-z0-9_.]{4,30}$'),
    CONSTRAINT valid_date_of_birth CHECK (date_of_birth <= CURRENT_DATE)
);

-- Index for username lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Index for email lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index for mobile number lookups
CREATE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile_number);