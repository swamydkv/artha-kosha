# Data Model: Infrastructure Modernization and Observability

**Feature**: 002-infra-observability-sessions  
**Date**: 2026-07-15  
**Status**: Draft

## Overview

This document defines the data model for the Infrastructure Modernization and Observability feature, including entities, relationships, validation rules, and database schema requirements. All data access uses PostgreSQL with pgx driver and handwritten SQL per constitutional requirements.

---

## Entity Definitions

### 1. Session

**Purpose**: Represents a user authentication session with secure cookie-based session management.

**Fields**:
- `id` (UUID, primary key) - Unique session identifier
- `user_id` (UUID, foreign key → users.id) - Associated user
- `created_at` (timestamptz, not null) - Session creation timestamp
- `last_activity_at` (timestamptz, not null) - Last user activity timestamp (sliding expiration)
- `expires_at` (timestamptz, not null) - Session expiration timestamp (absolute maximum)
- `revoked_at` (timestamptz, nullable) - Session revocation timestamp (null if active)
- `user_agent` (text, nullable) - Client user agent string
- `ip_address` (inet, nullable) - Client IP address
- `status` (session_status, not null) - Session status (active, expired, revoked)

**Validation Rules**:
- `id` must be a valid UUID
- `user_id` must reference an existing user
- `created_at` must be ≤ `last_activity_at`
- `last_activity_at` must be ≤ `expires_at`
- `expires_at` must be ≤ `created_at` + 8 hours (absolute maximum lifetime)
- `revoked_at` must be ≥ `created_at` if not null
- `status` must be one of: active, expired, revoked
- When `revoked_at` is set, `status` must be 'revoked'
- When `expires_at` < now(), `status` must be 'expired'

**State Transitions**:
```
[CREATED] → active
active → expired (automatic on expires_at < now())
active → revoked (user action or security event)
expired → (terminal state, no further transitions)
revoked → (terminal state, no further transitions)
```

**Indexes**:
- Primary key on `id`
- Index on `user_id` for user session queries
- Index on `expires_at` for cleanup jobs
- Composite index on `(user_id, status)` for active session lookups
- Index on `revoked_at` for audit queries

**Constraints**:
- Foreign key constraint on `user_id` → `users.id` (ON DELETE CASCADE)
- Check constraint: `expires_at <= created_at + interval '8 hours'`
- Check constraint: `last_activity_at >= created_at`
- Check constraint: `last_activity_at <= expires_at`

---

### 2. Audit Event

**Purpose**: Immutable record of business operations for security auditing and compliance.

**Fields**:
- `id` (UUID, primary key) - Unique audit event identifier
- `request_id` (text, not null) - Request correlation identifier
- `user_id` (UUID, nullable, foreign key → users.id) - User who performed the action
- `session_id` (UUID, nullable, foreign key → sessions.id) - Session in which action occurred
- `resource` (text, not null) - Resource type (e.g., "user", "account", "transaction")
- `resource_id` (text, nullable) - Specific resource identifier
- `action` (text, not null) - Action performed (e.g., "create", "update", "delete", "login", "logout")
- `result` (audit_result, not null) - Operation result (success, failure)
- `timestamp` (timestamptz, not null) - Event timestamp
- `user_agent` (text, nullable) - Client user agent string
- `client_ip` (inet, nullable) - Client IP address

**Validation Rules**:
- `id` must be a valid UUID
- `request_id` must not be empty
- `user_id` must reference an existing user if not null
- `session_id` must reference an existing session if not null
- `resource` must not be empty
- `action` must not be empty
- `result` must be one of: success, failure
- `timestamp` must be ≤ current time (no future events)
- If `user_id` is null, `session_id` must also be null (system events)

**State Transitions**:
- No state transitions - audit events are append-only immutable records

**Indexes**:
- Primary key on `id`
- Index on `request_id` for request correlation
- Index on `user_id` for user activity audits
- Index on `session_id` for session activity audits
- Index on `timestamp` for time-based queries
- Composite index on `(resource, resource_id)` for resource-specific audits
- Composite index on `(user_id, timestamp)` for user activity timelines

**Constraints**:
- Foreign key constraint on `user_id` → `users.id` (ON DELETE SET NULL)
- Foreign key constraint on `session_id` → `sessions.id` (ON DELETE SET NULL)
- Check constraint: `timestamp <= now()`
- No UPDATE or DELETE operations allowed (append-only)

---

### 3. Domain Event

**Purpose**: Represents a state change in the system for event-driven architecture and future integrations.

**Fields**:
- `id` (UUID, primary key) - Unique domain event identifier
- `event_type` (domain_event_type, not null) - Type of domain event
- `aggregate_id` (text, not null) - Aggregate root identifier (e.g., user ID, transaction ID)
- `aggregate_type` (text, not null) - Aggregate type (e.g., "user", "transaction", "account")
- `event_data` (jsonb, not null) - Event payload data
- `timestamp` (timestamptz, not null) - Event creation timestamp
- `processing_status` (processing_status, not null) - Processing status (pending, processed, failed)
- `retry_count` (integer, not null, default 0) - Number of processing retry attempts
- `processed_at` (timestamptz, nullable) - Timestamp when event was processed
- `error_message` (text, nullable) - Error message if processing failed

**Validation Rules**:
- `id` must be a valid UUID
- `event_type` must be one of: USER_REGISTERED, USER_LOGGED_IN, USER_LOGGED_OUT, ACCOUNT_CREATED, TRANSACTION_CREATED, BUDGET_CREATED
- `aggregate_id` must not be empty
- `aggregate_type` must not be empty
- `event_data` must be valid JSON
- `timestamp` must be ≤ current time
- `processing_status` must be one of: pending, processed, failed
- `retry_count` must be ≥ 0
- If `processing_status` = 'processed', `processed_at` must not be null
- If `processing_status` = 'failed', `error_message` must not be null

**State Transitions**:
```
[CREATED] → pending
pending → processed (successful processing)
pending → failed (processing error)
failed → pending (retry attempt)
processed → (terminal state)
```

**Indexes**:
- Primary key on `id`
- Index on `event_type` for event type queries
- Index on `aggregate_id` for aggregate-specific events
- Index on `processing_status` for pending event queries
- Index on `timestamp` for time-based queries
- Composite index on `(processing_status, timestamp)` for ordered processing
- Composite index on `(aggregate_type, aggregate_id)` for aggregate event streams

**Constraints**:
- Check constraint: `timestamp <= now()`
- Check constraint: `retry_count >= 0`
- Check constraint: `(processing_status = 'processed' AND processed_at IS NOT NULL) OR (processing_status != 'processed')`
- Check constraint: `(processing_status = 'failed' AND error_message IS NOT NULL) OR (processing_status != 'failed')`

---

### 4. Transactional Outbox

**Purpose**: Reliable event delivery mechanism ensuring domain events are persisted atomically with business operations.

**Fields**:
- `id` (UUID, primary key) - Unique outbox entry identifier
- `domain_event_id` (UUID, not null, foreign key → domain_events.id) - Reference to domain event
- `event_type` (text, not null) - Event type (denormalized for query performance)
- `payload` (jsonb, not null) - Event payload (denormalized for query performance)
- `created_at` (timestamptz, not null) - Outbox entry creation timestamp
- `processing_status` (processing_status, not null) - Processing status (pending, processed, failed)
- `retry_count` (integer, not null, default 0) - Number of delivery retry attempts
- `processed_at` (timestamptz, nullable) - Timestamp when event was delivered
- `error_message` (text, nullable) - Error message if delivery failed

**Validation Rules**:
- `id` must be a valid UUID
- `domain_event_id` must reference an existing domain event
- `event_type` must match the referenced domain event's type
- `payload` must match the referenced domain event's data
- `created_at` must be ≤ current time
- `processing_status` must be one of: pending, processed, failed
- `retry_count` must be ≥ 0
- If `processing_status` = 'processed', `processed_at` must not be null
- If `processing_status` = 'failed', `error_message` must not be null

**State Transitions**:
```
[CREATED] → pending (within same transaction as business operation)
pending → processed (successful delivery)
pending → failed (delivery error)
failed → pending (retry attempt)
processed → (terminal state)
```

**Indexes**:
- Primary key on `id`
- Index on `domain_event_id` for domain event reference
- Index on `processing_status` for pending event queries
- Index on `created_at` for time-based queries
- Composite index on `(processing_status, created_at)` for ordered delivery
- Partial index on `processing_status = 'pending'` for efficient polling

**Constraints**:
- Foreign key constraint on `domain_event_id` → `domain_events.id` (ON DELETE CASCADE)
- Check constraint: `created_at <= now()`
- Check constraint: `retry_count >= 0`
- Check constraint: `(processing_status = 'processed' AND processed_at IS NOT NULL) OR (processing_status != 'processed')`
- Check constraint: `(processing_status = 'failed' AND error_message IS NOT NULL) OR (processing_status != 'failed')`

---

## Enum Types

### session_status
```sql
CREATE TYPE session_status AS ENUM ('active', 'expired', 'revoked');
```

### audit_result
```sql
CREATE TYPE audit_result AS ENUM ('success', 'failure');
```

### domain_event_type
```sql
CREATE TYPE domain_event_type AS ENUM (
    'USER_REGISTERED',
    'USER_LOGGED_IN', 
    'USER_LOGGED_OUT',
    'ACCOUNT_CREATED',
    'TRANSACTION_CREATED',
    'BUDGET_CREATED'
);
```

### processing_status
```sql
CREATE TYPE processing_status AS ENUM ('pending', 'processed', 'failed');
```

---

## Entity Relationships

```
users (1) ←→ (N) sessions
       sessions (1) ←→ (N) audit_events
users (1) ←→ (N) audit_events
domain_events (1) ←→ (1) transactional_outbox
```

**Relationship Details**:
- A user can have multiple sessions (max 5 active per business rule)
- A session belongs to exactly one user
- An audit event can be associated with a session (optional)
- An audit event can be associated with a user (optional)
- A domain event has exactly one transactional outbox entry
- A transactional outbox entry references exactly one domain event

---

## Database Schema

### Table: sessions
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    user_agent TEXT,
    ip_address INET,
    status session_status NOT NULL DEFAULT 'active',
    
    CONSTRAINT check_expires_at CHECK (expires_at <= created_at + INTERVAL '8 hours'),
    CONSTRAINT check_activity_timing CHECK (last_activity_at >= created_at),
    CONSTRAINT check_activity_expires CHECK (last_activity_at <= expires_at)
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_user_status ON sessions(user_id, status);
CREATE INDEX idx_sessions_revoked_at ON sessions(revoked_at);
```

### Table: audit_events
```sql
CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL,
    resource TEXT NOT NULL,
    resource_id TEXT,
    action TEXT NOT NULL,
    result audit_result NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent TEXT,
    client_ip INET,
    
    CONSTRAINT check_timestamp CHECK (timestamp <= NOW())
);

CREATE INDEX idx_audit_events_request_id ON audit_events(request_id);
CREATE INDEX idx_audit_events_user_id ON audit_events(user_id);
CREATE INDEX idx_audit_events_session_id ON audit_events(session_id);
CREATE INDEX idx_audit_events_timestamp ON audit_events(timestamp);
CREATE INDEX idx_audit_events_resource ON audit_events(resource, resource_id);
CREATE INDEX idx_audit_events_user_time ON audit_events(user_id, timestamp);
```

### Table: domain_events
```sql
CREATE TABLE domain_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type domain_event_type NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    event_data JSONB NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status processing_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    
    CONSTRAINT check_timestamp CHECK (timestamp <= NOW()),
    CONSTRAINT check_retry_count CHECK (retry_count >= 0),
    CONSTRAINT check_processed CHECK (
        (processing_status = 'processed' AND processed_at IS NOT NULL) OR 
        (processing_status != 'processed')
    ),
    CONSTRAINT check_failed CHECK (
        (processing_status = 'failed' AND error_message IS NOT NULL) OR 
        (processing_status != 'failed')
    )
);

CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_id);
CREATE INDEX idx_domain_events_status ON domain_events(processing_status);
CREATE INDEX idx_domain_events_timestamp ON domain_events(timestamp);
CREATE INDEX idx_domain_events_status_time ON domain_events(processing_status, timestamp);
CREATE INDEX idx_domain_events_aggregate_stream ON domain_events(aggregate_type, aggregate_id);
```

### Table: transactional_outbox
```sql
CREATE TABLE transactional_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_event_id UUID NOT NULL REFERENCES domain_events(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processing_status processing_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    
    CONSTRAINT check_created_at CHECK (created_at <= NOW()),
    CONSTRAINT check_retry_count CHECK (retry_count >= 0),
    CONSTRAINT check_processed CHECK (
        (processing_status = 'processed' AND processed_at IS NOT NULL) OR 
        (processing_status != 'processed')
    ),
    CONSTRAINT check_failed CHECK (
        (processing_status = 'failed' AND error_message IS NOT NULL) OR 
        (processing_status != 'failed')
    )
);

CREATE INDEX idx_outbox_domain_event ON transactional_outbox(domain_event_id);
CREATE INDEX idx_outbox_status ON transactional_outbox(processing_status);
CREATE INDEX idx_outbox_created_at ON transactional_outbox(created_at);
CREATE INDEX idx_outbox_status_time ON transactional_outbox(processing_status, created_at);
CREATE INDEX idx_outbox_pending ON transactional_outbox(processing_status, created_at) 
    WHERE processing_status = 'pending';
```

---

## Data Consistency Rules

### ACID Transaction Requirements
Per FR-014 through FR-018, the following operations MUST execute within ACID transactions:

1. **User Registration**:
   - Insert user record
   - Create session
   - Generate audit event (USER_REGISTERED)
   - Generate domain event (USER_REGISTERED)
   - Insert outbox entry

2. **User Login**:
   - Validate credentials
   - Create session
   - Generate audit event (USER_LOGGED_IN)
   - Generate domain event (USER_LOGGED_IN)
   - Insert outbox entry

3. **User Logout**:
   - Revoke session
   - Generate audit event (USER_LOGGED_OUT)
   - Generate domain event (USER_LOGGED_OUT)
   - Insert outbox entry

4. **Account Creation**:
   - Insert account record
   - Generate audit event (ACCOUNT_CREATED)
   - Generate domain event (ACCOUNT_CREATED)
   - Insert outbox entry

5. **Transaction Creation**:
   - Insert transaction record
   - Update account balance
   - Generate audit event (TRANSACTION_CREATED)
   - Generate domain event (TRANSACTION_CREATED)
   - Insert outbox entry

6. **Budget Creation**:
   - Insert budget record
   - Generate audit event (BUDGET_CREATED)
   - Generate domain event (BUDGET_CREATED)
   - Insert outbox entry

### Transaction Failure Handling
- Any failure during transaction MUST result in complete rollback
- No partial data persistence allowed
- User receives clear error message
- No audit events or domain events persisted for failed operations

---

## Migration Strategy

### Migration Order
1. Create enum types (session_status, audit_result, domain_event_type, processing_status)
2. Create domain_events table
3. Create transactional_outbox table
4. Create sessions table
5. Create audit_events table
6. Create indexes
7. Create constraints

### Backward Compatibility
- Existing users table remains unchanged
- Existing authentication logic will be refactored to use sessions
- Migration will be additive (no destructive changes)
- Existing auth handlers will be updated incrementally

### Data Migration
- No existing session data to migrate (new feature)
- No existing audit data to migrate (new feature)
- No existing domain event data to migrate (new feature)
- Existing user accounts remain valid

---

## Security Considerations

### Sensitive Data Handling
- Passwords never stored in sessions (only user_id reference)
- IP addresses stored for audit trail (per constitutional requirements)
- User agent strings stored for audit trail
- No secrets or tokens in any entity
- Session IDs are cryptographically random UUIDs

### Access Control
- Session records accessible only to owning user
- Audit events accessible only to administrators (compliance)
- Domain events accessible only to internal consumers
- Outbox entries accessible only to internal event processors

### Retention Policies
- Sessions: Automatic cleanup on expiration + 30 days
- Audit events: Retain for 7 years (compliance requirement)
- Domain events: Retain for 1 year (business requirement)
- Outbox entries: Cleanup after successful processing + 7 days

---

## Performance Considerations

### Indexing Strategy
- All foreign keys indexed for join performance
- Time-based indexes for audit queries
- Status-based indexes for pending event processing
- Composite indexes for common query patterns

### Partitioning (Future)
- Consider partitioning audit_events by timestamp (monthly)
- Consider partitioning domain_events by timestamp (monthly)
- Not required for initial scale (<100 users)

### Cleanup Jobs
- Session cleanup: Run every 5 minutes, delete expired sessions > 30 days old
- Outbox cleanup: Run every hour, delete processed entries > 7 days old
- Domain event cleanup: Run daily, archive processed events > 1 year old

---

## Next Steps

1. Create SQL migration files in `apps/finance-api/migrations/`
2. Generate sqlc query definitions in `apps/finance-api/sql/`
3. Implement repository layer with type-safe SQL
4. Create integration tests for database operations
5. Implement transaction wrappers for ACID compliance
