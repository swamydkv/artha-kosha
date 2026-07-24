# Implementation Plan: Infrastructure Modernization and Observability

**Branch**: `002-infra-observability-sessions` | **Date**: 2026-07-15 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-infra-observability-sessions/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

This feature implements robust session management, authentication security, comprehensive observability, audit logging, domain events with transactional outbox, proper CORS configuration, and reliable database transactions for the ArthaKosha finance platform. The technical approach involves refactoring the existing authentication system to use persistent database sessions, implementing structured logging with correlation IDs, adding audit event tracking, implementing domain events with outbox pattern, and ensuring ACID transaction compliance for all financial operations.

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: pgx (PostgreSQL driver), log/slog (stdlib structured logging), go-chi/chi (HTTP router), go-chi/cors (CORS middleware), alexedwards/argon2id (password hashing)

**Storage**: PostgreSQL with versioned migrations

**Testing**: Go testing framework with contract, integration, and unit test layers

**Target Platform**: Linux server via Docker containers

**Project Type**: REST API web service (modular monolith architecture)

**Performance Goals**: 1000 concurrent authenticated users, under 2 second login session creation, <200ms p95 for API operations

**Constraints**: ACID transactions mandatory, no partial updates, structured logging with correlation IDs, HttpOnly session cookies only, no localStorage for auth, append-only audit logs

**Scale/Scope**: Single-family finance platform, modular monolith with domain-aligned modules, current user base <100 users

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Specification First
✅ Feature specification exists and is approved (spec.md)

### Test-Driven Development
✅ Feature specification includes acceptance criteria, will generate tasks before implementation

### OpenAPI-First Development
✅ OpenAPI specification updated with new session management endpoints (GET /session, DELETE /session, DELETE /sessions) and enhanced auth endpoints with session cookie headers and correlation ID headers

### Database-First Design
✅ Feature specification includes database schema requirements (sessions, audit_events, transactional_outbox tables)

### ACID Transactions
✅ Feature specification mandates ACID transactions for all financial operations

### Financial Data Integrity
✅ Feature specification maintains integer minor units for money (no changes needed)

### Modular Monolith Architecture
✅ Implementation follows modular monolith pattern within single finance-api application

### PostgreSQL-Only Data Access
✅ Uses PostgreSQL with pgx driver, no ORM frameworks

### Security and User Authorization
✅ Implements secure password hashing (Argon2id), HttpOnly cookies, audit logging
✅ Authentication abstraction layer implemented per FR-012 for future OIDC integration
✅ Custom authentication acceptable for current phase per constitution (OIDC planned for production)

### AI-Assisted Development
✅ AI assisting with specification and planning, human review required

### Repository Organization
✅ Follows established repository structure (apps/finance-api, specs/, infra/, docs/)

### Technology Standards
✅ Go 1.26+, PostgreSQL, REST API, Docker Compose

### Development Workflow
✅ Following specification-driven workflow with Mermaid diagrams in spec

### Session Management
✅ Implements persistent database sessions with HttpOnly cookies, sliding expiration, and session validation per constitutional requirements

### Observability and Audit Logging
✅ Implements structured logging with correlation IDs, audit events, and comprehensive observability per constitutional requirements

## Project Structure

### Documentation (this feature)

```text
specs/002-infra-observability-sessions/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
artha-kosha/
├── api/                      # OpenAPI specifications
├── apps/
│   ├── finance-api/          # Go backend API
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── auth/         # Authentication module
│   │   │   │   ├── contract/ # Contract tests
│   │   │   │   ├── integration/ # Integration tests
│   │   │   │   └── validation/ # Validation tests
│   │   │   ├── config/       # Configuration
│   │   │   ├── http/         # HTTP handlers and router
│   │   │   └── users/        # User domain module
│   │   ├── migrations/       # SQL migrations
│   │   └── sql/              # SQL query definitions
│   └── web/                  # Next.js frontend (future)
├── infra/
│   └── docker-compose.yml    # Infrastructure configuration
└── docs/                     # Architecture documentation
```

**Structure Decision**: This is a web application with Go backend (apps/finance-api) and Next.js frontend (apps/web). The backend follows a modular monolith architecture with domain-aligned modules (auth, users, etc.). All business logic resides in the finance-api application with clear separation between modules. The structure follows the constitutional requirement for modular monolith with PostgreSQL-only data access.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations requiring justification. All technical clarifications have been resolved in research.md and reflected in this plan.
