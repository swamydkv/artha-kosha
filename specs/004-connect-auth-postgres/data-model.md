# Data Model: Connect Auth System to PostgreSQL Users Table

## `users` Table
- `user_id` (UUID, Primary Key): Unique identifier.
- `username` (VARCHAR, Unique): Login name.
- `email` (VARCHAR, Unique): Contact email.
- `mobile_number` (VARCHAR, Unique): Phone number.
- `full_name` (VARCHAR): User's full name.
- `date_of_birth` (DATE): Date of birth.
- `password_hash` (VARCHAR): Secure Argon2id password hash.
- `created_at` (TIMESTAMPTZ): Registration timestamp.
- `updated_at` (TIMESTAMPTZ): Last modified timestamp.

## `sessions` Table
- `id` (UUID, Primary Key): Session token ID.
- `user_id` (UUID, Foreign Key -> `users.user_id`): Owner of the session.
- `expires_at` (TIMESTAMPTZ): Expiry time.

## `audit_events` Table
- `id` (UUID, Primary Key): Event identifier.
- `user_id` (UUID, Foreign Key -> `users.user_id`): User who performed the action.
- `action` (VARCHAR): Type of action (e.g., `user_registered`, `user_logged_in`).
- `result` (VARCHAR): Success/Failure.

## `domain_events` Table
- `id` (UUID, Primary Key): Event identifier.
- `aggregate_id` (VARCHAR): Usually matches `user_id`.
- `event_type` (VARCHAR): Type of event (e.g., `USER_REGISTERED`, `USER_LOGGED_IN`).
- `event_data` (JSONB): JSON payload of the event details.
