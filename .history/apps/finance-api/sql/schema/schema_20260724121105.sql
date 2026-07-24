-- Combined schema for sqlc generation (cleaned)

-- Enum types
CREATE TYPE session_status AS ENUM ('active','expired','revoked');
CREATE TYPE audit_result AS ENUM ('success','failure');
CREATE TYPE domain_event_type AS ENUM (
    'USER_REGISTERED', 'USER_LOGGED_IN', 'USER_LOGGED_OUT', 'ACCOUNT_CREATED', 'TRANSACTION_CREATED', 'BUDGET_CREATED'
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
