# Research Findings: Infrastructure Modernization and Observability

**Feature**: 002-infra-observability-sessions  
**Date**: 2026-07-15  
**Status**: Complete

## Overview

This document consolidates research findings for the technical decisions required to implement the Infrastructure Modernization and Observability feature. All clarifications from the Technical Context section of plan.md have been resolved.

---

## 1. Structured Logging Library

### Decision: **log/slog (Go 1.21+ standard library)**

### Rationale
- **Zero external dependencies** - Reduces supply chain risk for a finance platform
- **Go 1.26 available** - slog is mature and production-ready
- **Future-proof** - Ecosystem converging on slog; won't face deprecation
- **Performance sufficient** - ~101 ns/op with proper `LogAttrs()` usage
- **Native context support** - Built-in `InfoContext()`, `ErrorContext()` methods
- **Pluggable Handler architecture** - Can swap backends without changing application code
- **CORS headers already configured** - Infrastructure ready for correlation IDs

### Alternatives Considered
- **uber-go/zap**: Excellent performance (~51 ns/op) but adds external dependency
- **rs/zerolog**: Best raw performance (~25 ns/op) but external dependency with opinionated API
- **sirupsen/logrus**: **REJECTED** - Maintenance mode, poor performance (~9126 ns/op, 350x slower than zerolog)

### Implementation Pattern
- Custom `ContextHandler` for automatic field enrichment from context
- Middleware for request/correlation ID generation and propagation
- JSON output to stdout/stderr for platform-managed rotation
- Static fields (service, environment) added at logger initialization
- Dynamic fields (request_id, correlation_id, user_id, session_id) extracted from context

---

## 2. CORS Middleware

### Decision: **go-chi/cors** (with chi router)

### Rationale
- Designed specifically for chi router ecosystem
- All required features supported: configurable origins, methods, headers, credentials
- Supports wildcard "*" for headers (unlike gorilla/handlers)
- Has `AllowOriginFunc` for dynamic origin validation
- Returns appropriate 204 status for preflight (never 405 when configured correctly)
- Well-maintained fork of the popular rs/cors library

### Alternatives Considered
- **gin-contrib/cors**: Excellent but requires adopting Gin framework
- **gorilla/handlers CORS**: Does NOT support wildcard "*" for headers, must specify explicitly
- **rs/cors**: Framework-agnostic but can return 405 if using method-full patterns without proper setup
- **Custom implementation**: **REJECTED** - Error-prone, misses edge cases, no benefit over battle-tested libraries

### Configuration Requirements
- Origins: Configurable via application configuration (localhost:3000 for dev, production domains)
- Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS (OPTIONS critical to prevent 405)
- Headers: Content-Type, Authorization, X-Request-ID, X-Correlation-ID, X-Session-ID
- Exposed Headers: X-Request-ID, X-Correlation-ID
- Credentials: Enabled for HttpOnly cookie support
- MaxAge: 300-86400 seconds to reduce preflight requests

---

## 3. Password Hashing Library

### Decision: **alexedwards/argon2id**

### Rationale
- Uses `crypto/rand.Reader` for cryptographically secure salt generation
- Implements constant-time comparison using `crypto/subtle.ConstantTimeCompare`
- PHC-compliant output format: `$argon2id$v=19$m=65536,t=3,p=2$salt$hash`
- Enforces Argon2id variant only (prevents accidental use of weaker variants)
- Well-maintained, production-ready with clear API
- Inherits security from audited `golang.org/x/crypto/argon2` implementation

### Alternatives Considered
- **golang.org/x/crypto/argon2**: **REJECTED** - Raw primitive only, requires manual implementation of salt generation, encoding, and constant-time comparison
- **sixcolors/argon2id**: Good option for bcrypt migration (bcrypt-compatible API)
- **allisson/go-pwdhash**: Good for policy-based systems with automatic upgrade detection
- **matthewhartstonge/argon2**: Viable but constant-time comparison not explicitly documented

### OWASP Parameter Guidelines
Using RFC 9106 SECOND RECOMMENDED parameters for memory-constrained environments:
```go
params := &argon2id.Params{
    Memory:      65536,   // 64 MiB
    Iterations:  3,       // RFC SECOND RECOMMENDED
    Parallelism: 4,       // RFC recommendation
    SaltLength:  16,      // 128 bits
    KeyLength:   32,      // 256 bits
}
```

These parameters provide ~150-250ms per hash time with ~64 MiB RAM per hash, appropriate for a consumer finance application.

---

## 4. HTTP Router and Middleware Framework

### Decision: **go-chi/chi**

### Rationale
- **100% net/http compatible** - All handlers use standard `http.Handler` signature
- **Perfect for modular monolith architecture** - Route groups map naturally to modules/bounded contexts
- **Lightweight** - Core router is ~1000 LOC, no external dependencies
- **Excellent middleware ecosystem** - All middleware uses standard `func(http.Handler) http.Handler` signature
- **Context-based** - Built on Go's `context` package for request-scoped values
- **Active maintenance** - Latest release v5.3.0 (May 2026), recent security fixes
- **Production proven** - Used by Cloudflare, Heroku, 99Designs, Pressly
- **Easy migration** - Minimal changes from current stdlib `net/http` usage

### Alternatives Considered
- **gin-gonic/gin**: **REJECTED** - Fastest but not stdlib compatible, maintenance concerns (3-year release gaps), framework lock-in
- **labstack/echo**: Good middle ground but custom context type, more complex
- **gorilla/mux**: **REJECTED** - EOL suspected, slowest performance (13x slower than chi)
- **net/http (Go 1.22+)**: Viable but no built-in middleware chain, requires building custom utilities

### Middleware Ordering
Based on production best practices, the correct order is:
1. **Request ID** (outermost) - Generate correlation ID first
2. **Structured Logging** - Log requests with the request ID
3. **Recovery** - Catch panics from all inner middleware
4. **CORS/Security Headers** - Set early before auth checks
5. **Timeout** - Enforce timeout before expensive operations
6. **Authentication** - Validate identity
7. **Authorization** - Check permissions
8. **Router/Handler** (innermost)

---

## 5. OpenID Connect/Keycloak Integration Timeline

### Decision: **Phase 2 Implementation (Future Feature)**

### Rationale
- Current specification uses custom authentication with Argon2id password hashing
- Custom auth is acceptable for initial development phase
- OpenID Connect integration requires significant infrastructure setup (Keycloak deployment, configuration, user migration)
- Constitution specifies Keycloak for production but allows custom auth during development
- Current feature focuses on session management, observability, and infrastructure modernization
- OIDC integration should be a separate feature with dedicated specification and planning

### Implementation Approach
- Implement authentication abstraction layer as specified in FR-012
- Current implementation uses custom password-based authentication
- Future OIDC integration will plug into same abstraction layer
- Session management (current feature) works independently of authentication method
- Audit logging and observability work with both custom and OIDC authentication

---

## 6. OpenAPI Specification for Session Management

### Decision: **Create OpenAPI spec as part of Phase 1 design**

### Rationale
- Constitution requires OpenAPI-first development for all REST APIs
- New session management endpoints need specification: GET /session, DELETE /session, DELETE /sessions
- Enhanced auth endpoints need specification updates: login response to include session cookie
- OpenAPI spec will be created in `api/` directory following constitutional requirements
- Generated server types and validation will originate from OpenAPI specification
- **No API versioning required** - backend and frontend are developed together with no public API consumers
- Current APIs (/register, /login, /logout) will be retained and enhanced without versioning
- Breaking changes will be handled internally since both client and backend are under simultaneous development

### Endpoints to Specify
- `POST /login` - Update to include session cookie response (no versioning)
- `POST /logout` - Update to include session revocation (no versioning)
- `GET /session` - New endpoint for session validation (no versioning)
- `DELETE /session` - New endpoint for current session revocation (no versioning)
- `DELETE /sessions` - New endpoint for all sessions revocation (no versioning)

---

## Summary of Technical Decisions

| Component | Decision | Key Benefits |
|-----------|----------|--------------|
| **Structured Logging** | log/slog (stdlib) | Zero dependencies, future-proof, sufficient performance |
| **CORS Middleware** | go-chi/cors | Chi ecosystem integration, wildcard headers support, no 405 risk |
| **Password Hashing** | alexedwards/argon2id | Secure implementation, OWASP-compliant parameters, constant-time comparison |
| **HTTP Router** | go-chi/chi | Stdlib compatible, modular monolith design, active maintenance |
| **OIDC Integration** | Phase 2 (future) | Focus on current feature scope, proper abstraction layer |
| **OpenAPI Spec** | Phase 1 design | Constitutional compliance, contract-driven development |

All decisions align with the ArthaKosha constitution's principles:
- PostgreSQL-only data access ✅
- Modular monolith architecture ✅
- Security best practices ✅
- Zero unnecessary dependencies ✅
- Future-proof technology choices ✅

---

## Next Steps

With Phase 0 research complete, proceed to Phase 1 design:
1. Generate `data-model.md` with entities from feature spec
2. Define interface contracts in `contracts/` directory
3. Create `quickstart.md` validation guide
4. Re-evaluate Constitution Check post-design
