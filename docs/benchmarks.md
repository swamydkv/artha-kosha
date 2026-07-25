# Performance Benchmarks

This report tracks the historical progression of benchmark tests for various use cases across the backend, ensuring performance meets our targets as the application scales.

## Methodology & Running Benchmarks

Benchmarks are implemented as standard Go `testing.B` suites within the `integration` folders of their respective domain packages.

To run the full suite and update this file automatically:
```bash
cd apps/finance-api
go test -bench=. ./... | python3 ~/.gemini/antigravity-ide/brain/$(CONVERSATION_ID)/scratch/format_benchmarks.py
```
*(Note: A dedicated script or CI workflow will run this automatically on each feature completion).*

## Performance Targets

Based on the initial implementation plan:
- **Registration**: < 2 seconds
- **Login**: < 1 second
- **Logout**: < 500ms
- **Domain Operations (Accounts, Budgets, Tx)**: < 100ms per underlying service call.

## Performance Considerations

### In-Memory vs Database Persistence
Currently, benchmarking heavily reflects in-memory provider mocking to measure baseline domain logic overhead. 
- **Pros**: Lightning fast, precise measurement of business rules and algorithmic complexity (e.g. Bcrypt).
- **Cons**: Excludes Postgres latency.
*Note: Real PostgreSQL benchmarks will add ~10-50ms per operation depending on connection pooling (pgxpool) and index optimization.*

### Optimization Recommendations
1. **Database Strategy**: Implement caching (e.g., Redis) for frequent session checks, and maintain B-Tree indexes on core query parameters (username, email).
2. **Frontend Constraints**: Add debouncing and optimistic UI updates to mask network latency.

## Benchmark History

| Date | Component | Benchmark Name | Iterations | Time/Op | Notes |
|------|-----------|----------------|------------|---------|-------|
| 2026-07-25 09:44:19 | Accounts | `BenchmarkCreateAccount` | 49965052 | 20.04 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkRateLimiter_Allow` | 758779 | 1435 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkHashPassword` | 15 | 74584386 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkPasswordMatches` | 16 | 72362620 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkLogin` | 15 | 81319467 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkLogout` | 15 | 74317919 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkRegistration` | 100 | 37611986 ns/op | Baseline |
| 2026-07-25 09:44:19 | Auth | `BenchmarkRegistrationValidation` | 100 | 31264208 ns/op | Baseline |
| 2026-07-25 09:44:19 | Budgets | `BenchmarkCreateBudget` | 58345332 | 22.67 ns/op | Baseline |
| 2026-07-25 09:44:19 | Transactions | `BenchmarkCreateTransaction` | 51931372 | 21.39 ns/op | Baseline |
