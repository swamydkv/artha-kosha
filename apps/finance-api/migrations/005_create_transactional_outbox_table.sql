-- 005_create_transactional_outbox_table.sql
CREATE TABLE IF NOT EXISTS transactional_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_event_id UUID NOT NULL REFERENCES domain_events(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status processing_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_retry_at TIMESTAMPTZ,
    last_error TEXT,
    processed_at TIMESTAMPTZ,
    error_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_outbox_processing_status ON transactional_outbox(processing_status);
CREATE INDEX IF NOT EXISTS idx_outbox_created_at ON transactional_outbox(created_at);
-- Migration: 005_create_transactional_outbox_table.sql
CREATE TABLE IF NOT EXISTS transactional_outbox (
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

CREATE INDEX IF NOT EXISTS idx_outbox_domain_event ON transactional_outbox(domain_event_id);
CREATE INDEX IF NOT EXISTS idx_outbox_status ON transactional_outbox(processing_status);
