---
description: "Task list for Connect Auth System to PostgreSQL Users Table"
---

# Tasks: Connect Auth System to PostgreSQL Users Table

**Input**: Design documents from `/specs/004-connect-auth-postgres/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.
**Testing**: As per Constitution Principle II, Test-Driven Development is mandatory. Failing tests must be written first.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `apps/finance-api/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create initial SQL migration for users, sessions, audit_events, domain_events in apps/finance-api/migrations/00000X_init_auth.up.sql
- [x] T002 Configure sqlc schema and queries for auth data model in apps/finance-api/sql/users.sql
- [x] T003 Generate sqlc code via `sqlc generate` based on the new queries

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented
**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Implement database connection (pgxpool) initialization in apps/finance-api/cmd/main.go
- [x] T005 Implement background session cleanup worker in apps/finance-api/internal/sessions/worker.go
- [x] T006 Update cmd/main.go to start the background session cleanup worker

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Persist Registered Users (Priority: P1) 🎯 MVP

**Goal**: Persist newly registered users and sessions securely to PostgreSQL.
**Independent Test**: Register a new user, restart the application, and successfully log in using the same credentials.

### Tests for User Story 1 (MANDATORY TDD) ⚠️
> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T007 [P] [US1] Write failing unit tests for CreateUser and GetUserByUsername in apps/finance-api/internal/users/repository_sql_test.go
- [ ] T008 [P] [US1] Write failing integration test for user registration and login flows in apps/finance-api/internal/auth/provider_test.go

### Implementation for User Story 1

- [ ] T009 [US1] Implement CreateUser with atomic transactions for users, audit_events, and domain_events in apps/finance-api/internal/users/repository_sql.go
- [ ] T010 [US1] Implement GetUserByUsername in apps/finance-api/internal/users/repository_sql.go
- [ ] T011 [US1] Implement CreateSession with atomic transactions in apps/finance-api/internal/users/repository_sql.go
- [ ] T012 [US1] Refactor LocalAuthProvider in apps/finance-api/internal/auth/provider.go to use repository_sql for registration
- [ ] T013 [US1] Refactor LocalAuthProvider in apps/finance-api/internal/auth/provider.go to use repository_sql for login and session management

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Prevent Duplicate Accounts across Restarts (Priority: P2)

**Goal**: Reject new registrations that use an already registered username or email.
**Independent Test**: Register a user, restart the API, and attempt to register a new user with the same email/username to receive a duplicate error.

### Tests for User Story 2 (MANDATORY TDD) ⚠️
> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T014 [P] [US2] Write failing unit test for duplicate user registration in apps/finance-api/internal/users/repository_sql_test.go
- [ ] T015 [P] [US2] Write failing integration test for duplicate registration flow in apps/finance-api/internal/auth/provider_test.go

### Implementation for User Story 2

- [ ] T016 [US2] Update CreateUser method to handle and translate PostgreSQL unique constraint errors in apps/finance-api/internal/users/repository_sql.go
- [ ] T017 [US2] Handle duplicate account errors gracefully and return appropriate HTTP 409 Conflict in apps/finance-api/internal/auth/provider.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T018 [P] Verify exactly 100% statement coverage for backend packages and update docs/coverage.md
- [ ] T019 [P] Update docs/ architecture and call flow diagrams to reflect final implementation (Final Docs Rule)
- [ ] T020 [P] Verify structured logging rules (request IDs, correlation IDs) and audit logging for successful business operations

---

## Dependencies & Execution Order

### Phase Dependencies
- **Setup (Phase 1)**: Can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies
- **User Story 1 (P1)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 2 (P2)**: Depends on User Story 1 completion for full integration

### Parallel Opportunities
- TDD tests can be implemented concurrently with foundational tasks.
- SQL query creation (T002) and sqlc generation (T003) can be performed rapidly to unblock interface contracts.

---

## Implementation Strategy

### MVP First (User Story 1 Only)
1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently. Demo MVP.

### Incremental Delivery
1. Start User Story 2 (Duplicate Accounts prevention).
2. Validate independently.
3. Finish all polishing tasks and finalize documentation.
