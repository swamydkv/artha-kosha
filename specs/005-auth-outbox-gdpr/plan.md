# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]

**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

Fixing the UUID parsing bug for sessions, rewiring the `LoginService` to intercept auth domain events to the outbox, implementing GDPR-compliant soft deletion with an `archived_users` table, and introducing background sweepers to prune expired outbox and archive data.

## Technical Context

**Language/Version**: Go 1.26+, Next.js (TypeScript)

**Primary Dependencies**: `chi`, `sqlc`, `pgx`, `google/uuid`

**Storage**: PostgreSQL

**Testing**: Go `testing`

**Target Platform**: Linux / Docker

**Project Type**: REST API Web Service

**Performance Goals**: Negligible latency impact for authentication; outbox background workers must not block HTTP threads.

**Constraints**: Soft-deleted financial data MUST NOT break foreign keys (PII is anonymized, PK remains intact).

**Scale/Scope**: Designed for modular monolith; outbox cleanup ensures DB size remains manageable.

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
✅ Feature specification mandates ACID transactions for financial operations

### Financial Data Integrity
✅ Feature specification maintains integer minor units for money (if applicable)

### Modular Monolith Architecture
✅ Implementation follows modular monolith pattern within single application

### PostgreSQL-Only Data Access
✅ Uses PostgreSQL with pgx driver, no ORM frameworks

### Security and User Authorization
✅ Implements secure password hashing (Argon2id), proper session management, and authorization
✅ Authentication abstraction layer implemented for future OIDC integration
✅ Custom authentication acceptable for current phase per constitution (OIDC planned for production)

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
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
apps/finance-api/
├── cmd/main.go
├── internal/
│   ├── auth/
│   ├── users/
│   └── outbox/
├── migrations/
└── sql/
```

**Structure Decision**: Option 2: Web application / Modular monolith format as mandated by the constitution. Background workers will be spun up in `cmd/main.go` and logic housed in `internal/users` and `internal/outbox`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*No violations detected. Constitution check passes 100%.*
