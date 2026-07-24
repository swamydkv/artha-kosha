# Infrastructure Modernization Architecture

This document describes the updated architecture of the `artha-kosha` platform after implementing the infrastructure observability and transaction sessions features.

## Call Flow Diagram

```mermaid
sequenceDiagram
    participant Client
    participant Router as HTTP Router
    participant CORS as CORS Middleware
    participant Audit as Audit Middleware
    participant Auth as Auth Middleware
    participant Handler
    participant Service
    participant DB as Database

    Client->>Router: HTTP Request
    Router->>CORS: Preflight / Request
    CORS-->>Router: Allow Origin
    Router->>Audit: Record Start
    Router->>Auth: Validate Token
    Auth-->>Router: UserContext
    Router->>Handler: ServeHTTP
    Handler->>Service: Execute Logic
    Service->>DB: Begin Tx
    Service->>DB: Mutate State
    Service->>DB: Insert Audit Event
    Service->>DB: Commit Tx
    Service-->>Handler: Response
    Handler-->>Router: Response
    Router->>Audit: Record Finish
    Router-->>Client: HTTP Response
```

## Middleware Chain
1. Request ID Injection
2. Structured Logging (slog)
3. Panic Recovery
4. CORS Handling
5. Timeout (30s)
6. Audit Logging
7. Authentication
