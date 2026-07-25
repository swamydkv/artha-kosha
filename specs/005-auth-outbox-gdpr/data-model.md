# Data Model & Schema Additions

## New Entities

### `archived_users`
Stores the unscrubbed PII for soft-deleted accounts until the retention period expires.

**Fields**:
- `id` (UUID, Primary Key)
- `original_user_id` (UUID, Indexed) - Maps to the scrambled record in `users`
- `username` (VARCHAR)
- `email` (VARCHAR)
- `name` (VARCHAR)
- `mobile` (VARCHAR)
- `dob` (DATE)
- `password_hash` (VARCHAR)
- `archived_at` (TIMESTAMPTZ, Default NOW())
- `retention_expires_at` (TIMESTAMPTZ) - Used by background worker for cleanup

## Entity Modifications

### `users`
**Added Fields**:
- `deleted_at` (TIMESTAMPTZ, Nullable) - Indicates the account has been GDPR-deleted
- `is_archived` (BOOLEAN, Default FALSE) - Quick access flag to prevent logins for deleted accounts

*Note*: When an account is deleted, the existing PII fields (`email`, `name`, `mobile`, `username`) will be updated to deterministic or random anonymous hashes (e.g. `deleted_user_12345`) to preserve referential integrity for financial ledgers without retaining PII.

### `sessions` (Behavioral Change Only)
- No schema change, but the Go implementation MUST insert a true UUID for `id` rather than a `"session-xxx"` string.

### `transactional_outbox` (Behavioral Change Only)
- No schema change, but a background goroutine will issue `DELETE FROM transactional_outbox WHERE processing_status = 'processed' AND processed_at < NOW() - INTERVAL 'X days'`.
