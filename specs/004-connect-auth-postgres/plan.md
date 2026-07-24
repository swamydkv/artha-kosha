# Implementation Plan: Connect Auth System to PostgreSQL Users Table

**Branch**: `004-connect-auth-postgres` | **Date**: 2026-07-24 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/004-connect-auth-postgres/spec.md`

## Summary

The current `LocalAuthProvider` stores user registrations in an ephemeral map. This plan transitions user data to the PostgreSQL `users` table via `sqlc` generated queries, and introduces atomic transactions (`BEGIN`, `COMMIT`) for Registration and Login so that `domain_events` and `audit_events` are strictly synchronized with user and session creation.

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: `sqlc`, `pgx`

**Storage**: PostgreSQL

**Testing**: `go test`

**Target Platform**: Docker container running Linux/Alpine

**Project Type**: REST API

**Performance Goals**: <500ms p95 response time for registration/login

**Constraints**: Must strictly rollback on audit/domain event failure

**Scale/Scope**: Registration/login flow only

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Specification First
✅ Feature specification exists and is approved (spec.md)

### Test-Driven Development
✅ Feature specification includes acceptance criteria, will generate tasks before implementation

### OpenAPI-First Development
✅ OpenAPI specification is defined for all REST API endpoints

### Database-First Design
✅ Feature specification includes database schema requirements for persistent data

### ACID Transactions
✅ Feature specification mandates ACID transactions for financial operations (and auth operations)

### Financial Data Integrity
✅ N/A for Auth module

### Modular Monolith Architecture
✅ Implementation follows modular monolith pattern within single application

### PostgreSQL-Only Data Access
✅ Uses PostgreSQL with pgx driver, no ORM frameworks

### Security and User Authorization
✅ Implements secure password hashing (Argon2id), proper session management, and authorization
✅ Authentication abstraction layer implemented for future OIDC integration

### Session Management
✅ Implements persistent database sessions with HttpOnly cookies, sliding expiration, and session validation per constitutional requirements

### Observability and Audit Logging
✅ Implements structured logging with correlation IDs, audit events, and comprehensive observability per constitutional requirements

### AI-Assisted Development
✅ AI assisting with specification and planning, human review required

### Repository Organization
✅ Follows established repository structure (apps/, specs/, infra/, docs/)

### Technology Standards
✅ Follows constitutional technology standards (Go 1.26+, PostgreSQL, REST API, Docker Compose, etc.)

### Development Workflow
✅ Following specification-driven workflow with Mermaid diagrams in spec

## Project Structure

### Documentation (this feature)

```text
specs/004-connect-auth-postgres/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
apps/finance-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── auth/
│   │   └── provider.go
│   ├── sessions/
│   │   └── worker.go (new)
│   └── users/
│       └── repository_sql.go (new)
```

**Structure Decision**: The logic will be built directly into the existing `internal/users/` and `internal/auth/` directories. A new `repository_sql.go` will be added to `users` and the `provider.go` file will be refactored to use it. A background worker for session cleanup will be added.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
