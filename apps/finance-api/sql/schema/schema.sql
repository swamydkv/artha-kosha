-- Combined schema for sqlc generation (cleaned)

-- Enum types
CREATE TYPE session_status AS ENUM ('active','expired','revoked');
CREATE TYPE audit_result AS ENUM ('success','failure');
CREATE TYPE domain_event_type AS ENUM (
    'USER_REGISTERED', 'USER_LOGGED_IN', 'USER_LOGGED_OUT', 'USER_DELETED', 'ACCOUNT_CREATED', 'TRANSACTION_CREATED', 'BUDGET_CREATED'
);
CREATE TYPE processing_status AS ENUM ('pending','processed','failed');

-- Sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    user_agent TEXT,
    ip_address INET,
    status session_status NOT NULL DEFAULT 'active'
);

-- Users table (from migrations)
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    mobile_number TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Domain events
CREATE TABLE domain_events (
    id UUID PRIMARY KEY,
    event_type domain_event_type NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    event_data JSONB NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status processing_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    error_message TEXT
);

-- Transactional outbox
CREATE TABLE transactional_outbox (
    id UUID PRIMARY KEY,
    domain_event_id UUID NOT NULL REFERENCES domain_events(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status processing_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    last_retry_at TIMESTAMPTZ,
    last_error TEXT,
    error_message TEXT
);

-- Audit events
CREATE TABLE audit_events (
    id UUID PRIMARY KEY,
    request_id TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL,
    resource TEXT NOT NULL,
    resource_id TEXT,
    action TEXT NOT NULL,
    result audit_result NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent TEXT,
    client_ip INET
);

-- Accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    memo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Budgets table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
