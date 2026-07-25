# Phase 0: Outline & Research

## Technical Context Decisions

All technical context fields from the plan have been evaluated and resolved:

- **Language/Version**: Go 1.26+ for Backend, TypeScript/Next.js for Frontend
- **Primary Dependencies**: `chi` router, `pgx`/`sqlc` for Postgres, `google/uuid` for IDs
- **Storage**: PostgreSQL (via Docker)
- **Testing**: standard `go test` framework
- **Target Platform**: Docker-based Linux containers
- **Project Type**: Web Service (Finance API) + Web Client (Next.js)

## Research Tasks & Findings

### Decision 1: UUID Generation Strategy for Postgres
- **Decision**: Use `github.com/google/uuid` package to generate standard V4 UUIDs instead of the custom `"session-hex..."` string generator for Postgres sessions.
- **Rationale**: The `migrations/006_create_sessions_table.sql` explicitly types the `id` as `UUID`. Postgres throws a unique key violation if we attempt to insert the `uuid.Nil` fallback parsed from a malformed string.
- **Alternatives considered**: Changing the database schema to use `TEXT` for session IDs (rejected, violates existing DB design and indexing efficiency).

### Decision 2: GDPR Archiving Approach
- **Decision**: Introduce a new table `archived_users` rather than a separate database, and use a background worker goroutine in `cmd/main.go` to sweep expired records.
- **Rationale**: Keeping `archived_users` in the same Postgres DB simplifies the transaction boundaries for soft-deleting and moving PII atomically without requiring a distributed saga or two-phase commit.
- **Alternatives considered**: Shipping archived data to cold storage like AWS S3 (rejected, overkill for current modular monolith architecture).

### Decision 3: Outbox Polling vs Background Sweeping
- **Decision**: Augment the existing background worker (or add a dedicated ticker) to sweep `processed` outbox events older than `OUTBOX_RETENTION_DAYS`.
- **Rationale**: A simple time-based goroutine looping every hour/day is extremely lightweight and requires no external cron dependencies like Kubernetes CronJobs or Temporal.
- **Alternatives considered**: `pg_cron` extension in Postgres (rejected, adds infrastructure complexity to the DB tier).
