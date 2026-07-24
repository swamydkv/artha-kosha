# Implementation Tasks: 100% Coverage and Infrastructure Documentation

**Branch**: `[003-coverage-and-docs]`
**Spec**: [spec.md](file:///Users/swamydkv/Desktop/myProjects/artha-kosha/specs/003-coverage-and-docs/spec.md)
**Plan**: [plan.md](file:///Users/swamydkv/Desktop/myProjects/artha-kosha/specs/003-coverage-and-docs/plan.md)

## Phase 1: Implement Skipped Tests

- [x] T001: Implement `TestOutboxAtomicity` in `internal/events/integration/outbox_atomicity_test.go`
- [x] T002: Implement `TestSessionValidationLogic` (T032) in `internal/auth/validation/session_validation_test.go`
- [x] T003: Implement `TestSessionActiveExpiration` (T032A) in `internal/auth/integration/session_active_expiration_test.go`
- [x] T004: Implement `TestConcurrentLoginAttempts` (T032C) in `internal/auth/integration/concurrent_login_test.go`
- [x] T005: Implement `TestSessionLimitEnforcementNewDevice` (T032B) in `internal/auth/integration/session_limit_new_device_test.go`
- [x] T006: Implement `TestConcurrentSessionLimitEnforcement` (T031) in `internal/auth/integration/session_limit_test.go`
- [x] T007: Implement `TestSessionRevocation` (T030) in `internal/auth/integration/session_revocation_test.go`
- [x] T008: Implement `TestSessionExpirationHandling` (T029) in `internal/auth/integration/session_expiration_test.go`
- [x] T009: Implement `TestSessionCreationAndValidation` (T028) in `internal/auth/integration/session_flow_test.go`
- [x] T010: Implement `TestGetSessionContract` (T025) in `internal/auth/contract/session_contract_test.go`
- [x] T011: Implement `TestDeleteSessionsContract` (T027) in `internal/auth/contract/sessions_revoke_all_contract_test.go`
- [x] T012: Implement `TestDeleteSessionContract` (T026) in `internal/auth/contract/session_revoke_contract_test.go`

## Phase 2: Achieve 100% Statement Coverage

- [x] T013: Write tests for `cmd` and `cmd/migrate`
- [x] T014: Write tests for `internal/config`
- [x] T015: Write tests for `internal/database`
- [x] T016: Write tests for `internal/domain`
- [x] T017: Write tests for `internal/sessions` (e.g. `InMemoryRepo`, `PostgresRepo`, `Service`)
- [x] T018: Write tests for `internal/users`
- [x] T019: Write tests for `internal/audit`
- [x] T020: Write tests for `internal/sqlc` (generated models/queries)
- [x] T021: Add remaining edge-case tests to `internal/accounts`, `internal/transactions`, and `internal/budgets`
- [x] T022: Add remaining edge-case tests to `internal/auth` and `internal/http` packages
- [x] T023: Verify exactly 100.0% coverage across all packages using `go test -cover ./...`

## Phase 3: Documentation and Constitution Updates

- [x] T024: Update docs to include a Mermaid ER diagram of the database tables.
- [x] T025: Append 100% test coverage principle to `.specify/memory/constitution.md`
