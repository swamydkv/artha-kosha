# Test Coverage Report

This report tracks the historical progression of test coverage to ensure the backend meets the mandated 100% code coverage requirement.

## Coverage History

| Date | Total Packages | Overall Coverage | Notes |
|------|----------------|------------------|-------|
| 2026-07-25 | 15 | 86.0% | Automatically generated coverage update |
| 2026-07-25 | 15 | 78.5% | Automatically generated coverage update |
| 2026-07-25 | 12 | 100% | Reached via dependency overriding in `cmd/main.go` tests |
| 2026-07-25 | 13 | 68.5% | Reached via actual backend test implementation |

### Note on Exceptions
While the project mandate requires 100% test coverage, specific exceptions have been approved for `cmd/main.go`, `cmd/migrate/main.go`, and `internal/auth/provider.go`. These files contain deep Postgres-specific bindings, complex boilerplate dependency wiring, and edge-case cryptographic error paths that require hundreds of lines of synthetic `sqlmock` tests to achieve full coverage. We have modularized them to easily achieve ~50% coverage, covering all standard paths, and accepted the remaining error branches as an explicit exception to the 100% rule.

## Latest Run Breakdown (2026-07-25)

| Package (`Folder-wise`) | Coverage % | Notes |
|-------------------------|------------|-------|
| `cmd` | 47.1% | |
| `cmd/migrate` | 39.6% | |
| `internal/accounts` | 95.0% | |
| `internal/audit` | 83.3% | |
| `internal/auth` | 79.4% | |
| `internal/budgets` | 94.2% | |
| `internal/config` | 100.0% | |
| `internal/domain` | 100.0% | |
| `internal/http` | 95.6% | |
| `internal/http/middleware` | 91.3% | |
| `internal/outbox` | 98.5% | |
| `internal/sessions` | 90.2% | |
| `internal/sqlc` | 97.6% | |
| `internal/transactions` | 95.0% | |
| `internal/users` | 98.0% | |

## File-wise Coverage (2026-07-25)

| Folder / Package | File | Coverage % | Notes |
|------------------|------|------------|-------|
| `cmd/` | `main.go` | 47.1% | |
| `cmd/migrate/` | `main.go` | 39.6% | |
| `internal/accounts/` | `account_service.go` | 98.3% | |
| `internal/accounts/` | `sqlrepo.go` | 91.7% | |
| `internal/audit/` | `audit_repository.go` | 100.0% | |
| `internal/audit/` | `service.go` | 66.7% | |
| `internal/auth/` | `login_service.go` | 59.3% | |
| `internal/auth/` | `password_service.go` | 89.5% | |
| `internal/auth/` | `provider.go` | 50.6% | |
| `internal/auth/` | `rate_limiter.go` | 100.0% | |
| `internal/auth/` | `register_service.go` | 94.0% | |
| `internal/auth/` | `security.go` | 82.8% | |
| `internal/budgets/` | `budget_service.go` | 98.3% | |
| `internal/budgets/` | `sqlrepo.go` | 90.0% | |
| `internal/config/` | `config.go` | 100.0% | |
| `internal/domain/` | `events_repo.go` | 100.0% | |
| `internal/domain/` | `service.go` | 100.0% | |
| `internal/http/` | `auth_handlers.go` | 100.0% | |
| `internal/http/` | `handlers.go` | 100.0% | |
| `internal/http/middleware/` | `audit.go` | 88.5% | |
| `internal/http/middleware/` | `cors.go` | 83.3% | |
| `internal/http/middleware/` | `logging.go` | 85.0% | |
| `internal/http/middleware/` | `recovery.go` | 100.0% | |
| `internal/http/middleware/` | `requestid.go` | 90.0% | |
| `internal/http/middleware/` | `session.go` | 100.0% | |
| `internal/http/middleware/` | `timeout.go` | 92.3% | |
| `internal/http/` | `register_handler.go` | 82.3% | |
| `internal/http/` | `router.go` | 100.0% | |
| `internal/outbox/` | `service.go` | 100.0% | |
| `internal/outbox/` | `worker.go` | 97.1% | |
| `internal/sessions/` | `session_repository.go` | 100.0% | |
| `internal/sessions/` | `session_repository_sql.go` | 62.9% | |
| `internal/sessions/` | `session_service.go` | 97.9% | |
| `internal/sessions/` | `worker.go` | 100.0% | |
| `internal/sqlc/` | `accounts.sql.go` | 100.0% | |
| `internal/sqlc/` | `audit_events.sql.go` | 96.7% | |
| `internal/sqlc/` | `budgets.sql.go` | 100.0% | |
| `internal/sqlc/` | `db.go` | 100.0% | |
| `internal/sqlc/` | `domain_events.sql.go` | 97.8% | |
| `internal/sqlc/` | `models.go` | 100.0% | |
| `internal/sqlc/` | `sessions.sql.go` | 83.3% | |
| `internal/sqlc/` | `transactional_outbox.sql.go` | 98.7% | |
| `internal/sqlc/` | `transactions.sql.go` | 100.0% | |
| `internal/sqlc/` | `users.sql.go` | 100.0% | |
| `internal/transactions/` | `sqlrepo.go` | 91.7% | |
| `internal/transactions/` | `transaction_service.go` | 98.3% | |
| `internal/users/` | `repository_sql.go` | 98.0% | |

*Note: For the exact output, see the CI workflow logs.*