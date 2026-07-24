# Database Schema

This document outlines the database schema for the Artha Kosha Finance API.

## Entity-Relationship Diagram

```mermaid
erDiagram
    users {
        UUID user_id PK
        TEXT full_name
        DATE date_of_birth
        TEXT mobile_number UK
        TEXT email UK
        TEXT username UK
        TEXT password_hash
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    sessions {
        UUID id PK
        UUID user_id FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ last_activity_at
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ revoked_at
        TEXT user_agent
        INET ip_address
        session_status status
    }

    domain_events {
        UUID id PK
        domain_event_type event_type
        TEXT aggregate_id
        TEXT aggregate_type
        JSONB event_data
        TIMESTAMPTZ timestamp
        processing_status processing_status
        INTEGER retry_count
        TIMESTAMPTZ processed_at
        TEXT error_message
    }

    transactional_outbox {
        UUID id PK
        UUID domain_event_id FK
        TEXT event_type
        JSONB payload
        TIMESTAMPTZ created_at
        processing_status processing_status
        INTEGER retry_count
        TIMESTAMPTZ processed_at
        TEXT error_message
    }

    audit_events {
        UUID id PK
        TEXT request_id
        UUID user_id FK
        UUID session_id FK
        TEXT resource
        TEXT resource_id
        TEXT action
        audit_result result
        TIMESTAMPTZ timestamp
        TEXT user_agent
        INET client_ip
    }

    users ||--o{ sessions : "has"
    users ||--o{ audit_events : "generates"
    sessions ||--o{ audit_events : "associated_with"
    domain_events ||--o{ transactional_outbox : "produces"
```

## Enum Types

- `session_status`: `active`, `expired`, `revoked`
- `domain_event_type`: `UserRegistered`, `UserLoggedIn`, `UserLoggedOut`, etc.
- `processing_status`: `pending`, `processing`, `completed`, `failed`
- `audit_result`: `success`, `failure`
