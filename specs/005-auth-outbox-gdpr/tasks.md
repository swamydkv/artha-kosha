---

description: "Task list for auth-outbox-gdpr implementation"
---

# Tasks: auth-outbox-gdpr

**Input**: Design documents from `/specs/005-auth-outbox-gdpr/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure
*(No shared setup tasks required for this feature extension on existing infrastructure.)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T001 Create database migration for `archived_users` and `users` changes in `migrations/007_gdpr_archiving.sql`
- [X] T002 Add `OutboxRetentionDays` and `ArchiveRetentionDays` to `internal/config/config.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Secure Session Resumption (Priority: P1) 🎯 MVP

**Goal**: Fix the UUID session parsing bug that prevents subsequent logins.

**Independent Test**: Register, logout, and login again successfully.

### Implementation for User Story 1

- [X] T003 [US1] Fix UUID generation for Postgres sessions in `internal/auth/provider.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 3 - GDPR Data Archiving & Deletion (Priority: P1)

**Goal**: Implement the GDPR-compliant account deletion and PII anonymization.

**Independent Test**: Trigger `DELETE /user/account` and confirm PII is scrubbed in `users` and moved to `archived_users`.

### Implementation for User Story 3

- [X] T004 [P] [US3] Add SQL models and repository methods for `ArchiveUser` and `DeleteUser` in `internal/users/repository_sql.go`, including session revocation.
- [X] T005 [P] [US3] Create the HTTP handler for `DELETE /user/account` and wire it in `cmd/main.go`
- [X] T005a [P] [US3] Frontend: Implement "Danger Zone" account deletion UI in Profile/Settings page.
- [X] T005b [P] [US3] Frontend: Add "DELETE" confirmation text prompt before submitting the deletion request.

**Checkpoint**: Account deletion works independently.

---

## Phase 5: User Story 2 - Accurate Audit and Event Tracking (Priority: P2)

**Goal**: Wire the outbox domain events and audit events during registration, login, and logout.

**Independent Test**: Perform auth operations and query the outbox tables to verify insertions.

### Implementation for User Story 2

- [X] T006 [P] [US2] Update `internal/users/repository_sql.go` to explicitly write to `transactional_outbox` during `CreateUser` transaction.
- [X] T007 [P] [US2] Rewire `cmd/main.go` so the HTTP router points to `LoginService` wrapper instead of directly to `LocalAuthProvider`.

**Checkpoint**: Outbox events and audit events correctly emit.

---

## Phase 6: User Story 4 - Background Data Lifecycle Management (Priority: P3)

**Goal**: Create sweeping workers to prevent table bloat.

**Independent Test**: Wait for worker tick and observe old entries being deleted.

### Implementation for User Story 4

- [X] T008 [P] [US4] Add `DeleteProcessed` method to `internal/outbox/service.go` and its repository.
- [X] T009 [P] [US4] Add `PruneArchivedUsers` method to `internal/users/repository_sql.go`.
- [X] T010 [US4] Launch goroutine sweepers in `cmd/main.go` utilizing the config retention days.

**Checkpoint**: All user stories should now be independently functional.

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T011 Run `quickstart.md` validation.
- [X] T012 Write comprehensive unit tests for `cmd/main.go` to achieve 100% backend test coverage.
- [X] T013 Generate and update `docs/coverage.md` with the 100% coverage report.
- [X] T014 Finalize architecture and write documentation to the `docs/` directory.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: N/A
- **Foundational (Phase 2)**: BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - Can proceed in priority order (P1 → P1 → P2 → P3) or parallel.
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (Session Bug)**: Independent.
- **US3 (GDPR Archiving)**: Independent.
- **US2 (Outbox Wiring)**: Independent.
- **US4 (Workers)**: Dependent on US2 (for outbox generation) and US3 (for archive generation) to be manually verifiable.
