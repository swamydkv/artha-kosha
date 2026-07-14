# Data Model: User Registration and Login

## Entities

### User

| Field | Type | Constraints | Notes |
|---|---|---|---|
| user_id | UUID v7 | Primary key, not null | Generated automatically |
| full_name | text | Not null | Display name for onboarding |
| date_of_birth | date | Not null | Must not be in the future |
| mobile_number | text | Not null, unique | Required for MVP |
| email | text | Not null, unique | Must be valid email |
| username | text | Not null, unique | 4–30 chars; letters, numbers, `_`, `.` |
| password_hash | text | Not null | One-way hash only |
| created_at | timestamptz | Not null | System timestamp |
| updated_at | timestamptz | Not null | System timestamp |

### Session

| Field | Type | Constraints | Notes |
|---|---|---|---|
| session_id | UUID v7 | Primary key | Generated on login |
| user_id | UUID v7 | Foreign key to `users.user_id` | Owner of the session |
| created_at | timestamptz | Not null | Session creation time |
| expires_at | timestamptz | Not null | Session expiration time |
| revoked_at | timestamptz | Nullable | Set on logout |

## Relationships

- One `User` may have many `Session` records over time.
- Each `Session` belongs to exactly one `User`.

## Validation Rules

- Username must match the required alphanumeric and punctuation rule.
- Password must satisfy the minimum complexity requirements.
- Email must be syntactically valid.
- Date of birth must be in the past or present, never in the future.
- Mobile number must be unique if required by the MVP acceptance flow.

## Transactional Guarantees

- Account registration is persisted atomically in one transaction.
- Session creation and invalidation are transactional where supported.
