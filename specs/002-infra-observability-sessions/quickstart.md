# Quickstart Validation Guide: Infrastructure Modernization and Observability

**Feature**: 002-infra-observability-sessions  
**Date**: 2026-07-15  
**Status**: Draft

## Overview

This guide provides runnable validation scenarios to verify that the Infrastructure Modernization and Observability feature works end-to-end. Use these scenarios to validate session management, authentication security, observability, audit logging, and domain events.

---

## Prerequisites

### Infrastructure Setup
```bash
# Start PostgreSQL and other infrastructure
cd /Users/swamydkv/Desktop/myProjects/artha-kosha/infra
docker-compose up -d

# Verify PostgreSQL is running
docker-compose ps

# Run database migrations
cd ../apps/finance-api
go run cmd/main.go migrate up
```

### Application Build
```bash
# Build the application
cd /Users/swamydkv/Desktop/myProjects/artha-kosha/apps/finance-api
go build -o finance-api cmd/main.go

# Or run directly
go run cmd/main.go
```

### Environment Configuration
Ensure the following environment variables are set:
```bash
DATABASE_URL=postgres://user:password@localhost:5432/artha_kosha
SESSION_TIMEOUT_MINUTES=15
SESSION_MAX_HOURS=8
MAX_CONCURRENT_SESSIONS=5
CORS_ORIGINS=http://localhost:3000
LOG_LEVEL=info
```

---

## Validation Scenarios

### Scenario 1: User Registration with Session Creation

**Objective**: Verify user registration creates a session with HttpOnly cookie and generates audit/domain events.

**Steps**:
```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "date_of_birth": "1990-01-01",
    "mobile_number": "+1234567890",
    "email": "test@example.com",
    "username": "testuser",
    "password": "SecurePassword123!",
    "confirm_password": "SecurePassword123!"
  }' \
  -v -c cookies.txt
```

**Expected Results**:
- HTTP 201 response
- Response contains `user_id` and `username`
- `Set-Cookie` header with `session_id` (HttpOnly, Path=/, Max-Age=28800, SameSite=Strict)
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- Database: User record created in `users` table
- Database: Session record created in `sessions` table with status 'active'
- Database: Audit event created in `audit_events` table (action: USER_REGISTERED)
- Database: Domain event created in `domain_events` table (type: USER_REGISTERED)
- Database: Outbox entry created in `transactional_outbox` table

**Validation Queries**:
```sql
-- Check user exists
SELECT id, username FROM users WHERE username = 'testuser';

-- Check session exists and is active
SELECT id, user_id, status, expires_at 
FROM sessions 
WHERE user_id = (SELECT id FROM users WHERE username = 'testuser')
AND status = 'active';

-- Check audit event
SELECT * FROM audit_events 
WHERE action = 'USER_REGISTERED' 
ORDER BY timestamp DESC LIMIT 1;

-- Check domain event
SELECT * FROM domain_events 
WHERE event_type = 'USER_REGISTERED' 
ORDER BY timestamp DESC LIMIT 1;

-- Check outbox entry
SELECT * FROM transactional_outbox 
WHERE processing_status = 'pending' 
ORDER BY created_at DESC LIMIT 1;
```

---

### Scenario 2: User Login with Session Management

**Objective**: Verify user login creates/renews session with proper cookie handling.

**Steps**:
```bash
# Login with existing user
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePassword123!"
  }' \
  -v -c cookies.txt
```

**Expected Results**:
- HTTP 200 response
- Response contains `user_id`, `session_id`, and `welcome_message`
- `Set-Cookie` header with `session_id` (HttpOnly, Path=/, Max-Age=28800, SameSite=Strict)
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- Database: Session record created or updated in `sessions` table
- Database: `last_activity_at` updated to current timestamp
- Database: Audit event created (action: USER_LOGGED_IN)
- Database: Domain event created (type: USER_LOGGED_IN)
- Database: Outbox entry created

**Validation Queries**:
```sql
-- Check session is active
SELECT id, user_id, last_activity_at, expires_at, status 
FROM sessions 
WHERE user_id = (SELECT id FROM users WHERE username = 'testuser')
AND status = 'active'
ORDER BY created_at DESC LIMIT 1;

-- Check audit event
SELECT * FROM audit_events 
WHERE action = 'USER_LOGGED_IN' 
ORDER BY timestamp DESC LIMIT 1;

-- Check domain event
SELECT * FROM domain_events 
WHERE event_type = 'USER_LOGGED_IN' 
ORDER BY timestamp DESC LIMIT 1;
```

---

### Scenario 3: Session Validation and Persistence

**Objective**: Verify session persists across requests and sliding expiration works.

**Steps**:
```bash
# Validate current session
curl -X GET http://localhost:8080/session \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -v
```

**Expected Results**:
- HTTP 200 response
- Response contains `user_id`, `session_id`, `created_at`, `last_activity_at`, `expires_at`, `status`
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- Database: `last_activity_at` updated to current timestamp (sliding expiration)
- Status is 'active'

**Validation Query**:
```sql
-- Check session activity was updated
SELECT id, last_activity_at, expires_at, status 
FROM sessions 
WHERE id = '<session_id_from_response>';
```

---

### Scenario 4: Session Expiration Handling

**Objective**: Verify expired sessions are rejected and handled properly.

**Steps**:
```bash
# Manually expire a session in database
UPDATE sessions 
SET expires_at = NOW() - INTERVAL '1 hour',
    status = 'expired'
WHERE user_id = (SELECT id FROM users WHERE username = 'testuser')
AND status = 'active';

# Attempt to use expired session
curl -X GET http://localhost:8080/session \
  -b cookies.txt \
  -v
```

**Expected Results**:
- HTTP 401 response (Unauthorized)
- Error message indicates session expired
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- No database updates (session remains expired)

---

### Scenario 5: Session Revocation (Single Device)

**Objective**: Verify single session revocation works correctly.

**Steps**:
```bash
# Revoke current session
curl -X DELETE http://localhost:8080/session \
  -b cookies.txt \
  -v -c cookies.txt
```

**Expected Results**:
- HTTP 200 response
- `Set-Cookie` header clears session cookie (Max-Age=0)
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- Database: Session status updated to 'revoked'
- Database: `revoked_at` set to current timestamp
- Database: Audit event created (action: USER_LOGGED_OUT)
- Database: Domain event created (type: USER_LOGGED_OUT)

**Validation Queries**:
```sql
-- Check session is revoked
SELECT id, status, revoked_at 
FROM sessions 
WHERE id = '<session_id>';

-- Check audit event
SELECT * FROM audit_events 
WHERE action = 'USER_LOGGED_OUT' 
ORDER BY timestamp DESC LIMIT 1;

-- Check domain event
SELECT * FROM domain_events 
WHERE event_type = 'USER_LOGGED_OUT' 
ORDER BY timestamp DESC LIMIT 1;
```

---

### Scenario 6: Session Revocation (All Devices)

**Objective**: Verify all sessions for a user can be revoked.

**Steps**:
```bash
# Create multiple sessions (simulate multiple devices)
# First login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "SecurePassword123!"}' \
  -c cookies1.txt

# Second login (different device)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "SecurePassword123!"}' \
  -c cookies2.txt

# Revoke all sessions
curl -X DELETE http://localhost:8080/sessions \
  -b cookies1.txt \
  -v
```

**Expected Results**:
- HTTP 200 response
- `Set-Cookie` header clears session cookie
- `X-Request-ID` header present
- `X-Correlation-ID` header present
- Database: ALL sessions for user have status 'revoked'
- Database: ALL sessions have `revoked_at` set
- Database: Audit event created (action: USER_LOGGED_OUT)
- Database: Domain event created (type: USER_LOGGED_OUT)

**Validation Query**:
```sql
-- Check all sessions are revoked
SELECT id, status, revoked_at 
FROM sessions 
WHERE user_id = (SELECT id FROM users WHERE username = 'testuser');
```

---

### Scenario 7: Concurrent Session Limit

**Objective**: Verify maximum 5 concurrent sessions per user is enforced.

**Steps**:
```bash
# Create 5 sessions (should succeed)
for i in {1..5}; do
  curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser", "password": "SecurePassword123!"}' \
    -c cookies$i.txt
done

# Attempt 6th session (should fail)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "SecurePassword123!"}' \
  -v
```

**Expected Results**:
- First 5 login attempts: HTTP 200 (success)
- 6th login attempt: HTTP 403 (Forbidden) or appropriate error
- Error message indicates session limit reached
- Database: Only 5 active sessions exist for user

**Validation Query**:
```sql
-- Check active session count
SELECT COUNT(*) 
FROM sessions 
WHERE user_id = (SELECT id FROM users WHERE username = 'testuser')
AND status = 'active';
```

---

### Scenario 8: Structured Logging with Correlation IDs

**Objective**: Verify structured logging includes all required fields and correlation IDs propagate.

**Steps**:
```bash
# Make a request with custom correlation ID
curl -X GET http://localhost:8080/session \
  -H "X-Correlation-ID: test-correlation-123" \
  -H "X-Request-ID: test-request-456" \
  -b cookies.txt \
  -v
```

**Expected Results**:
- HTTP 200 response
- `X-Request-ID` header present (either provided or generated)
- `X-Correlation-ID` header present (echoed back)
- Application logs contain:
  - `timestamp`
  - `level` (info, error, etc.)
  - `service` (artha-kosha-finance-api)
  - `component` (auth, session, etc.)
  - `operation` (session_validate, etc.)
  - `request_id` (matches header)
  - `correlation_id` (matches header)
  - `user_id` (from session)
  - `session_id` (from cookie)
  - `duration` (request processing time)
  - `message` (descriptive log message)

**Log Validation**:
```bash
# Check application logs
docker-compose logs finance-api | grep test-correlation-123
```

Expected log format (JSON):
```json
{
  "timestamp": "2026-07-15T10:30:00Z",
  "level": "info",
  "service": "artha-kosha-finance-api",
  "component": "session",
  "operation": "session_validate",
  "request_id": "test-request-456",
  "correlation_id": "test-correlation-123",
  "user_id": "uuid",
  "session_id": "uuid",
  "duration": "15ms",
  "message": "Session validated successfully"
}
```

---

### Scenario 9: CORS Configuration

**Objective**: Verify CORS headers are properly configured and preflight requests work.

**Steps**:
```bash
# Preflight OPTIONS request
curl -X OPTIONS http://localhost:8080/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  -v

# Actual POST request with CORS
curl -X POST http://localhost:8080/login \
  -H "Origin: http://localhost:3000" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "SecurePassword123!"}' \
  -v
```

**Expected Results**:
- OPTIONS request: HTTP 200 (never 405)
- CORS headers present:
  - `Access-Control-Allow-Origin: http://localhost:3000`
  - `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
  - `Access-Control-Allow-Headers: Content-Type, Authorization, X-Request-ID, X-Correlation-ID, X-Session-ID`
  - `Access-Control-Expose-Headers: X-Request-ID, X-Correlation-ID`
  - `Access-Control-Allow-Credentials: true`
  - `Access-Control-Max-Age: 300`
- POST request: HTTP 200 with same CORS headers
- Unauthorized origins rejected (CORS policy violation)

---

### Scenario 10: ACID Transaction Compliance

**Objective**: Verify database operations are atomic and rollback on failure.

**Steps**:
```bash
# Attempt transaction that will fail (e.g., invalid data)
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "account_id": "invalid-uuid",
    "amount": -100,
    "type": "expense"
  }' \
  -v
```

**Expected Results**:
- HTTP 400 or 422 response (validation error)
- No partial data persisted in database
- No audit event created for failed operation
- No domain event created for failed operation
- No outbox entry created for failed operation
- Database state unchanged (atomic rollback)

**Validation Queries**:
```sql
-- Verify no transaction was created
SELECT COUNT(*) FROM transactions WHERE amount = -100;

-- Verify no audit event for failed operation
SELECT COUNT(*) FROM audit_events 
WHERE resource = 'transaction' 
AND action = 'create' 
AND result = 'failure';

-- Verify no domain event for failed operation
SELECT COUNT(*) FROM domain_events 
WHERE event_type = 'TRANSACTION_CREATED';
```

---

### Scenario 11: Password Security

**Objective**: Verify passwords are hashed using Argon2id with proper parameters.

**Steps**:
```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Security Test",
    "date_of_birth": "1990-01-01",
    "mobile_number": "+1234567890",
    "email": "security@example.com",
    "username": "secuser",
    "password": "AnotherSecurePassword456!",
    "confirm_password": "AnotherSecurePassword456!"
  }'
```

**Expected Results**:
- HTTP 201 response
- Database: Password hash is in PHC format: `$argon2id$v=19$m=65536,t=3,p=2$...`
- Database: Hash parameters match OWASP recommendations (m=65536, t=3, p=2)
- Database: Salt is 16 bytes (128 bits)
- Database: Hash length is 32 bytes (256 bits)
- Plain text password NOT stored anywhere

**Validation Query**:
```sql
-- Check password hash format
SELECT password_hash 
FROM users 
WHERE username = 'secuser';
```

Expected hash format: `$argon2id$v=19$m=65536,t=3,p=2$<base64_salt>$<base64_hash>`

---

### Scenario 12: Domain Event Processing

**Objective**: Verify domain events are stored in outbox and can be processed.

**Steps**:
```bash
# Create a transaction (generates domain event)
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "account_id": "<valid-account-id>",
    "amount": 5000,
    "type": "expense",
    "category": "groceries",
    "description": "Test transaction"
  }'
```

**Expected Results**:
- HTTP 201 response
- Database: Domain event created in `domain_events` table
- Database: Outbox entry created in `transactional_outbox` table
- Database: Outbox entry has `processing_status = 'pending'`
- Domain event and outbox entry are in same transaction (atomic)

**Validation Queries**:
```sql
-- Check domain event
SELECT * FROM domain_events 
WHERE event_type = 'TRANSACTION_CREATED' 
ORDER BY timestamp DESC LIMIT 1;

-- Check outbox entry
SELECT * FROM transactional_outbox 
WHERE processing_status = 'pending' 
ORDER BY created_at DESC LIMIT 1;

-- Verify atomicity (both should exist or neither)
SELECT COUNT(*) FROM domain_events de
JOIN transactional_outbox to ON de.id = to.domain_event_id
WHERE de.event_type = 'TRANSACTION_CREATED';
```

---

## Cleanup

```bash
# Stop infrastructure
cd /Users/swamydkv/Desktop/myProjects/artha-kosha/infra
docker-compose down

# Clean test data (optional)
psql -U user -d artha_kosha -c "DELETE FROM users WHERE username IN ('testuser', 'secuser');"
```

---

## Success Criteria

All scenarios are successful when:
- ✅ HTTP responses match expected status codes
- ✅ Response headers include `X-Request-ID` and `X-Correlation-ID`
- ✅ Session cookies are HttpOnly with correct attributes
- ✅ Database state matches expected results
- ✅ Audit events are created for all operations
- ✅ Domain events are created for state changes
- ✅ Outbox entries are created atomically with operations
- ✅ Structured logs contain all required fields
- ✅ CORS headers are properly configured
- ✅ Transactions are atomic (rollback on failure)
- ✅ Passwords are hashed with Argon2id
- ✅ Session limits are enforced

---

## References

- [API Contract](./contracts/auth-api.yaml)
- [Data Model](./data-model.md)
- [Feature Specification](./spec.md)
- [Research Findings](./research.md)
