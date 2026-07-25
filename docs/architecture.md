# ArthaKosha Architecture

## Overview
ArthaKosha is a private personal and family finance management platform. It follows a Modular Monolith Architecture as mandated by its Constitution.

## Technology Stack
- **Frontend**: Next.js (TypeScript) + TailwindCSS
- **Backend**: Go 1.26+ (REST API)
- **Database**: PostgreSQL (via pgx and sqlc)
- **Authentication**: Custom session-based authentication with Argon2id (OpenID Connect via Keycloak planned)

## System Design
The backend is structured into cohesive domain modules:
- `accounts`: Manages financial accounts.
- `transactions`: Manages ledgers and records.
- `budgets`: Manages budgetary rules.
- `users`: Manages user lifecycle, GDPR soft-deletions, and archiving.
- `auth`: Manages sessions and security.
- `audit`: Immutable audit logging.
- `outbox`: Transactional outbox for event dispatching.

### GDPR Compliance & Soft-Deletions
When an account deletion is requested:
1. Active sessions are revoked.
2. PII is scrubbed (anonymized) in the main `users` table to preserve financial ledgers.
3. Raw PII is moved to `archived_users` with an expiration timestamp.
4. Background workers periodically prune expired records from `archived_users`.

### Observability
- All state-changing operations are wrapped in transactions and produce immutable audit events.
- Domain events are written to the `transactional_outbox` for asynchronous processing.
- Background workers prune old, processed outbox events to prevent database bloat.
