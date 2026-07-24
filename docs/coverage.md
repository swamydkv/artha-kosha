# Code Coverage Report

## Repo-Wise Coverage

| Total Statements | Covered Statements | Coverage % |
| --- | --- | --- |
| 1168 | 870 | 74.5% |

## Folder-Wise Coverage

| Folder | Statements | Covered | Coverage % |
| --- | --- | --- | --- |
| `artha-kosha/apps/finance-api/cmd` | 42 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/cmd/migrate` | 22 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/internal/accounts` | 34 | 23 | 67.6% |
| `artha-kosha/apps/finance-api/internal/audit` | 21 | 20 | 95.2% |
| `artha-kosha/apps/finance-api/internal/auth` | 279 | 229 | 82.1% |
| `artha-kosha/apps/finance-api/internal/budgets` | 33 | 22 | 66.7% |
| `artha-kosha/apps/finance-api/internal/config` | 4 | 4 | 100.0% |
| `artha-kosha/apps/finance-api/internal/database` | 10 | 10 | 100.0% |
| `artha-kosha/apps/finance-api/internal/domain` | 35 | 35 | 100.0% |
| `artha-kosha/apps/finance-api/internal/http` | 141 | 61 | 43.3% |
| `artha-kosha/apps/finance-api/internal/http/middleware` | 102 | 76 | 74.5% |
| `artha-kosha/apps/finance-api/internal/outbox` | 72 | 53 | 73.6% |
| `artha-kosha/apps/finance-api/internal/sessions` | 113 | 96 | 85.0% |
| `artha-kosha/apps/finance-api/internal/sqlc` | 165 | 160 | 97.0% |
| `artha-kosha/apps/finance-api/internal/transactions` | 34 | 23 | 67.6% |
| `artha-kosha/apps/finance-api/internal/users` | 61 | 58 | 95.1% |

## File-Wise Coverage

| File | Statements | Covered | Coverage % |
| --- | --- | --- | --- |
| `artha-kosha/apps/finance-api/cmd/main.go` | 42 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/cmd/migrate/main.go` | 22 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/internal/accounts/account_service.go` | 27 | 17 | 63.0% |
| `artha-kosha/apps/finance-api/internal/accounts/sqlrepo.go` | 7 | 6 | 85.7% |
| `artha-kosha/apps/finance-api/internal/audit/audit_repository.go` | 12 | 12 | 100.0% |
| `artha-kosha/apps/finance-api/internal/audit/service.go` | 9 | 8 | 88.9% |
| `artha-kosha/apps/finance-api/internal/auth/login_service.go` | 51 | 34 | 66.7% |
| `artha-kosha/apps/finance-api/internal/auth/password_service.go` | 19 | 17 | 89.5% |
| `artha-kosha/apps/finance-api/internal/auth/provider.go` | 133 | 119 | 89.5% |
| `artha-kosha/apps/finance-api/internal/auth/rate_limiter.go` | 15 | 15 | 100.0% |
| `artha-kosha/apps/finance-api/internal/auth/register_service.go` | 49 | 34 | 69.4% |
| `artha-kosha/apps/finance-api/internal/auth/security.go` | 12 | 10 | 83.3% |
| `artha-kosha/apps/finance-api/internal/budgets/budget_service.go` | 27 | 17 | 63.0% |
| `artha-kosha/apps/finance-api/internal/budgets/sqlrepo.go` | 6 | 5 | 83.3% |
| `artha-kosha/apps/finance-api/internal/config/config.go` | 4 | 4 | 100.0% |
| `artha-kosha/apps/finance-api/internal/database/transaction.go` | 10 | 10 | 100.0% |
| `artha-kosha/apps/finance-api/internal/domain/events_repo.go` | 19 | 19 | 100.0% |
| `artha-kosha/apps/finance-api/internal/domain/service.go` | 16 | 16 | 100.0% |
| `artha-kosha/apps/finance-api/internal/http/auth_handlers.go` | 63 | 13 | 20.6% |
| `artha-kosha/apps/finance-api/internal/http/handlers.go` | 39 | 27 | 69.2% |
| `artha-kosha/apps/finance-api/internal/http/middleware/audit.go` | 16 | 15 | 93.8% |
| `artha-kosha/apps/finance-api/internal/http/middleware/cors.go` | 6 | 6 | 100.0% |
| `artha-kosha/apps/finance-api/internal/http/middleware/logging.go` | 25 | 20 | 80.0% |
| `artha-kosha/apps/finance-api/internal/http/middleware/recovery.go` | 7 | 7 | 100.0% |
| `artha-kosha/apps/finance-api/internal/http/middleware/requestid.go` | 23 | 15 | 65.2% |
| `artha-kosha/apps/finance-api/internal/http/middleware/session.go` | 12 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/internal/http/middleware/timeout.go` | 13 | 13 | 100.0% |
| `artha-kosha/apps/finance-api/internal/http/register_handler.go` | 18 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/internal/http/router.go` | 21 | 21 | 100.0% |
| `artha-kosha/apps/finance-api/internal/outbox/service.go` | 33 | 32 | 97.0% |
| `artha-kosha/apps/finance-api/internal/outbox/worker.go` | 39 | 21 | 53.8% |
| `artha-kosha/apps/finance-api/internal/sessions/session_repository.go` | 43 | 43 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sessions/session_repository_sql.go` | 43 | 37 | 86.0% |
| `artha-kosha/apps/finance-api/internal/sessions/session_service.go` | 17 | 16 | 94.1% |
| `artha-kosha/apps/finance-api/internal/sessions/worker.go` | 10 | 0 | 0.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/accounts.sql.go` | 2 | 2 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/audit_events.sql.go` | 17 | 16 | 94.1% |
| `artha-kosha/apps/finance-api/internal/sqlc/budgets.sql.go` | 2 | 2 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/db.go` | 2 | 2 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/domain_events.sql.go` | 19 | 18 | 94.7% |
| `artha-kosha/apps/finance-api/internal/sqlc/models.go` | 52 | 52 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/sessions.sql.go` | 14 | 12 | 85.7% |
| `artha-kosha/apps/finance-api/internal/sqlc/transactional_outbox.sql.go` | 23 | 22 | 95.7% |
| `artha-kosha/apps/finance-api/internal/sqlc/transactions.sql.go` | 2 | 2 | 100.0% |
| `artha-kosha/apps/finance-api/internal/sqlc/users.sql.go` | 32 | 32 | 100.0% |
| `artha-kosha/apps/finance-api/internal/transactions/sqlrepo.go` | 7 | 6 | 85.7% |
| `artha-kosha/apps/finance-api/internal/transactions/transaction_service.go` | 27 | 17 | 63.0% |
| `artha-kosha/apps/finance-api/internal/users/repository_sql.go` | 61 | 58 | 95.1% |
