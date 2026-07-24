<!--
Sync Impact Report
- Version change: 1.3.0 → 1.4.0
- Modified principles: added session management requirements, observability and audit logging requirements, updated transactional outbox from optional to standard pattern, clarified session-based authentication during development
- Added sections: Principle XII (Session Management), Principle XIII (Observability and Audit Logging)
- Removed sections: none
- Templates requiring updates: spec-template.md ✅ updated to include session management and observability requirements, plan-template.md ✅ updated to include session management and observability in constitution check
- Follow-up TODOs: none
-->

# ArthaKosha Constitution

ArthaKosha is a private personal and family finance management platform for tracking accounts, income, expenses, budgets, goals, receipts, reports, and future financial insights. It is not a banking system, payment gateway, or money transfer platform. The application manages financial records and metadata while preserving correctness, auditability, and user privacy.

## Core Principles

### I. Specification First
The approved specification is the source of truth. Every feature MUST follow the lifecycle: Specification → Clarification → Technical Plan → Tasks → Tests → Implementation → Review → Specification Conformance. Implementation MUST NOT diverge from the approved specification.

### II. Test-Driven Development
Test-driven development is mandatory. Every feature MUST define acceptance criteria before implementation, create failing tests before production code, include unit and integration coverage, add regression tests for every defect fix, and pass all automated tests before merge. No production code MAY exist without corresponding automated tests.

### III. OpenAPI-First Development
All REST APIs MUST be designed using OpenAPI before implementation. The OpenAPI specification is the API contract. Generated server types, request validation, and client SDKs MUST originate from the OpenAPI specification whenever practical.

### IV. Database-First Design
Every feature that affects persistent data MUST define its domain model, database schema, SQL migration plan, constraints, foreign keys, and indexes. Database changes MUST be introduced only through versioned SQL migrations.

### V. ACID Transactions
Financial operations MUST preserve consistency. Related database changes MUST execute within PostgreSQL ACID transactions. Partial updates are prohibited. Rollback MUST occur whenever any operation in the transaction fails.

### VI. Financial Data Integrity
Money MUST always be stored as integer minor units. Floating-point types MUST NEVER be used for monetary values. Currency MUST always be explicit. Historical financial records MUST remain traceable. Corrections MUST be represented through explicit adjustment or reversal operations rather than silently modifying historical records.

### VII. Modular Monolith Architecture
The system MUST begin as a modular monolith. The backend MUST consist of clearly separated business modules within a single deployable application. Microservices, distributed messaging, caching infrastructure, Kubernetes, or event streaming platforms MUST NOT be introduced until justified by measurable operational requirements. A transactional outbox MUST be implemented for state-changing operations to support audit logging, reporting, notifications, categorization, forecasting, and AI capabilities.

### VIII. PostgreSQL-Only Data Access
PostgreSQL is the system of record. Data access MUST use pgx, sqlc, and handwritten SQL. ORM frameworks MUST NOT be used.

### IX. Security and User Authorization
Authentication MUST use OpenID Connect, with Keycloak planned for production. During development, session-based authentication with secure password hashing (Argon2id) is acceptable. Authorization MUST always be based on the authenticated user. Every operation MUST verify that the authenticated user has permission to access or modify the requested resource. Secrets MUST be managed using HashiCorp Vault. Secrets, credentials, and production data MUST NEVER be committed to source control. Only synthetic data MAY be used in automated tests and AI-assisted development.

### X. AI-Assisted Development
AI tools MAY assist development but MUST NOT replace engineering judgment. AI MAY generate code, tests, SQL, documentation, and boilerplate. Developers remain responsible for reviewing and approving specifications, architecture, database schema, security decisions, API contracts, and generated code. Generated code MUST NEVER override approved specifications.

### XI. Repository Organization
The repository structure is an architectural decision and MUST remain consistent throughout the project. AI-generated code and developer contributions MUST conform to the established repository layout. Changes to the repository organization MUST undergo architectural review and MUST be documented through an Architecture Decision Record (ADR).

The project MUST use the following top-level structure:

```text
artha-kosha/
├── AGENTS.md
├── README.md
├── specs/
├── api/
├── apps/
│   ├── finance-api/
│   └── web/
├── infra/
└── docs/
```

The responsibilities of each directory are:

- `specs/` – Specifications, technical plans, implementation tasks, and acceptance criteria.
- `api/` – OpenAPI specifications and API contracts.
- `apps/finance-api/` – Go backend application source code.
- `apps/web/` – Next.js frontend application source code.
- `infra/` – Infrastructure configuration, Docker Compose files, deployment assets, and environment configuration.
- `docs/` – Architecture Decision Records (ADRs), domain glossary, and supporting documentation.

Within the backend application:

- Versioned SQL migrations MUST reside in a dedicated `migrations` directory.
- SQL query definitions for sqlc MUST reside in a dedicated `sql` directory.
- Business functionality MUST be organized into cohesive modules that align with the domain model.
- Shared infrastructure and reusable components MUST be separated from domain-specific business logic.

AI-generated code MUST reuse existing modules whenever possible and MUST NOT introduce new directory structures, duplicate architectural patterns, or reorganize the repository without explicit approval.

Feature specifications MUST describe business behavior and requirements, not implementation file paths or package names. File placement and project organization MUST follow this constitution and the project's architectural conventions.

### XII. Session Management
User sessions MUST be managed securely with persistent database storage. Sessions MUST use HttpOnly cookies for web clients and MUST NOT use localStorage for authentication state. Sessions MUST implement sliding expiration with configurable inactivity timeout and absolute maximum lifetime. Sessions MUST support concurrent session limits per user and provide explicit revocation capabilities (single device and all devices). Session validation MUST occur on every authenticated request.

### XIII. Observability and Audit Logging
The system MUST implement comprehensive structured logging with consistent log levels and formats. Every log entry MUST include timestamp, log level, service, component, operation, request ID, correlation ID, user ID, session ID, duration, and message. The system MUST assign unique request IDs and correlation IDs to every request and propagate them through all operations. The system MUST generate immutable audit events for every successful business operation with complete context (user, action, resource, timestamp, user agent, client IP). Audit events MUST be append-only and never modified. Sensitive data (passwords, tokens, secrets, hashes) MUST NEVER appear in log files.

## Technology Standards

- Backend: Go 1.26+
- API: REST
- API Specification: OpenAPI
- Database: PostgreSQL
- SQL Access: pgx + sqlc
- Web: Next.js + TypeScript
- Mobile (future): React Native
- Infrastructure: Docker Compose
- Secrets: HashiCorp Vault
- Version Control: Git
- Logging: log/slog (Go 1.21+ standard library)
- Password Hashing: Argon2id (alexedwards/argon2id)
- HTTP Router: go-chi/chi
- CORS: go-chi/cors

## Development Workflow & Definition of Done

A feature is complete only when all of the following are true:

- The specification is approved.
- Clarifications are resolved.
- The system design is documented using Mermaid diagrams.
- User flows and call flows are explicitly captured in the approved specification using Mermaid diagrams.
- The OpenAPI contract is finalized.
- The database schema and migrations are reviewed.
- Session management requirements are implemented (HttpOnly cookies, sliding expiration, validation on every request).
- Observability requirements are implemented (structured logging, correlation IDs, audit events).
- Security requirements are implemented (Argon2id password hashing, CORS configuration).
- Transactional outbox is implemented for state-changing operations.
- Acceptance tests pass.
- Unit tests pass.
- Integration tests pass.
- The implementation conforms to the approved specification.
- Documentation is updated in the `docs/` folder as the final step of the feature lifecycle, reflecting the finalized architecture, user flow, call flow, and supporting design rationale.

The normal delivery workflow MUST be specification-driven, contract-driven, test-driven, migration-driven, and documentation-finalized. Feature specifications MUST contain the exact user flow and call flow for the feature using Mermaid diagrams. All architecture, user-flow, and workflow artifacts MUST be represented in Mermaid. Documentation updates in the `docs/` folder MUST occur after implementation details have stabilized and MUST reflect the final architecture, not draft or exploratory structure.

## Governance

This constitution supersedes ad hoc development practices for ArthaKosha. Amendment proposals MUST document the change, explain the governance impact, and identify any migration or workflow implications. Any change that modifies a principle, adds a mandatory rule, or tightens a compliance gate MUST be reviewed before it is adopted.

All pull requests and design reviews MUST verify compliance with this constitution. Any exception MUST be justified in writing, approved by the responsible technical owner, and limited to the smallest possible scope. Complexity or architecture changes that violate a principle require explicit rationale and evidence that a simpler alternative was evaluated and rejected.

The constitution uses semantic versioning: MAJOR for backward-incompatible governance or principle removals, MINOR for new principles or materially expanded guidance, and PATCH for clarifications, wording fixes, and non-semantic refinements. The version line MUST reflect both the current release and the adoption dates.

**Version**: 1.4.0 | **Ratified**: 2026-07-13 | **Last Amended**: 2026-07-24
