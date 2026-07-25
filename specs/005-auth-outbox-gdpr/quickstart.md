# Quickstart & Validation Guide

## Prerequisites
- Docker Compose running the Postgres database
- The `finance-api` backend running locally (`go run cmd/main.go`)

## Validation 1: Proper Session UUID (Bug Fix)
1. Register a new user:
   ```bash
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```
2. Log in (creates session 1):
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```
3. Log out.
4. Log in **again** (creates session 2). 
   - *Expected Outcome*: It succeeds and returns a 200 OK without Postgres Primary Key violations.

## Validation 2: Outbox & Audit Events
1. Connect to Postgres:
   ```bash
   docker exec -it artha-kosha-postgres-1 psql -U postgres -d artha_kosha
   ```
2. Query the tables:
   ```sql
   SELECT * FROM domain_events WHERE event_type IN ('USER_REGISTERED', 'USER_LOGGED_IN', 'USER_LOGGED_OUT');
   SELECT * FROM audit_events;
   SELECT * FROM transactional_outbox;
   ```
   - *Expected Outcome*: You should see matching entries for all the operations you just performed.

## Validation 3: GDPR Deletion
1. Initiate the deletion (using the auth token/cookie from your active session):
   ```bash
   curl -X DELETE http://localhost:8080/user/account \
     -H "Cookie: session_id=YOUR_SESSION_UUID" \
     -H "Content-Type: application/json" \
     -d '{"confirmation": "DELETE"}'
   ```
2. Check the database:
   ```sql
   SELECT username, is_archived, deleted_at FROM users WHERE id = '...';
   SELECT * FROM archived_users WHERE original_user_id = '...';
   ```
   - *Expected Outcome*: The main `users` table has an anonymized username and `is_archived = true`. The `archived_users` table contains the original unscrubbed PII and a `retention_expires_at` date 30 days in the future.
