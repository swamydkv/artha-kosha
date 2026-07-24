# Research & Decisions: Connect Auth System to PostgreSQL Users Table

## 1. Transactional Registration & Login
**Decision:** Implement a `CreateUserTx` and `LoginTx` method in `users.Repository` and `sessions.Repository` (or handle it in a higher-level Service/Unit of Work) to wrap the `users`/`sessions`, `domain_events`, and `audit_events` inserts in a single transaction.
**Rationale:** The Functional Requirements mandate strict transaction safety (fail and rollback on audit/domain event failure). Since `sqlc` generates queries using `context.Context` and `*sql.DB`/`*sql.Tx`, we can use `db.BeginTx()` to start a transaction, use `q.WithTx(tx)` to run `CreateUser`, `CreateDomainEvent`, and `CreateAuditEvent` inside the transaction, and then `tx.Commit()`.
**Alternatives considered:**
- Event sourcing pattern (too complex for the current architecture).
- Two-phase commit (unnecessary since all tables are in the same PostgreSQL database).

## 2. Session Cleanup Background Worker
**Decision:** Create a goroutine that runs on a configured interval (e.g., every 1 hour) in `apps/finance-api/cmd/main.go` or `internal/sessions/worker.go` that executes a `DELETE FROM sessions WHERE expires_at < NOW() - INTERVAL '30 days'` query.
**Rationale:** FR-009b requires a configurable background cleanup process to physically delete expired sessions after 30 days. Go routines are lightweight and perfectly suited for this periodic cleanup.
**Alternatives considered:**
- PostgreSQL `pg_cron` extension (requires modifying DB server configuration, not easily portable).
- Deleting synchronously during `GetSession` (does not clean up sessions that are never accessed again).

## 3. Wiring Auth Provider to `users` Repository
**Decision:** Modify `internal/auth/provider.go` to remove `users map[string]*localUser`. Inject a `users.Repository` interface implementation (`users.SQLRepository`). Refactor `Register` and `Login` to invoke the `users.Repository` methods.
**Rationale:** This fulfills the primary goal of the feature. We already have `sqlc` generated queries in `internal/sqlc/users.sql.go` for creating and fetching users.
**Alternatives considered:**
- Continuing to use the in-memory map (rejected, defeats the purpose of the persistence layer).
