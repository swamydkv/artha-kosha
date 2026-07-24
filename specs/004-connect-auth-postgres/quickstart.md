# Quickstart: Connect Auth System to PostgreSQL Users Table

## Overview
This guide demonstrates how to validate that the Auth system is successfully integrated with the PostgreSQL users table and that transactions (including `audit_events` and `domain_events`) are processed correctly.

## Prerequisites
- The `artha-kosha-postgres-1` container must be running.
- The `artha-kosha-finance-api-1` container must be running.
- Access to `curl` or a REST client.

## Validation Scenarios

### 1. Register a New User
Run the following `curl` command to register a new user:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "date_of_birth": "1990-01-01",
    "mobile_number": "1234567890",
    "email": "test@example.com",
    "username": "testuser",
    "password": "Password123!",
    "confirm_password": "Password123!"
  }'
```
**Expected Outcome:** 
- A `201 Created` response.

### 2. Verify Persistence in PostgreSQL
Open a psql session and check the tables:
```bash
docker exec -it artha-kosha-postgres-1 psql -U postgres -d artha_kosha
```
```sql
SELECT * FROM users WHERE username = 'testuser';
SELECT * FROM domain_events WHERE event_type = 'USER_REGISTERED';
SELECT * FROM audit_events WHERE action = 'user_registered';
```
**Expected Outcome:** 
- You should see exactly 1 row returned for each query, with matching identifiers.

### 3. Verify Persistence Across Restarts
Restart the API container to simulate a deployment or crash:
```bash
docker restart artha-kosha-finance-api-1
```
Then, try to log in:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Password123!"
  }'
```
**Expected Outcome:**
- A `200 OK` response with a valid session ID.
- The previous in-memory implementation would have failed with "invalid credentials" after a restart, proving that PostgreSQL persistence is now active!
