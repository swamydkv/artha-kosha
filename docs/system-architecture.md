# ArthaKosha System Architecture

**Version**: 1.1.0  
**Status**: Active  

## Overview
This document represents the latest architecture, system design, and call flows of the ArthaKosha backend (`finance-api`). The application follows a Modular Monolith architecture, using Go, PostgreSQL, and follows strict adherence to ACID transactions and the Transactional Outbox pattern for domain events.

For database entity relationships, please refer to the [Database ER Diagram](database.md).

---

## 1. High-Level Component Architecture

```mermaid
flowchart TB
    subgraph "Frontend Layer"
        WEB[Next.js Web Client]
        REG[Registration Page]
        LOG[Login Page]
        HOME[Authenticated App]
    end
    
    subgraph "API Layer (HTTP / Middleware)"
        API[Go REST API & Router]
        MW_CORS[CORS Middleware]
        MW_LOG[Structured Logging]
        MW_AUDIT[Audit Middleware]
        MW_AUTH[Session Auth Middleware]
        
        API --> MW_CORS
        MW_CORS --> MW_LOG
        MW_LOG --> MW_AUDIT
        MW_AUDIT --> MW_AUTH
    end
    
    subgraph "Domain Layer (Business Logic)"
        AUTH_SVC[Auth Service]
        SESS_SVC[Session Service]
        DOMAIN_SVC[Domain Events Service]
        OUTBOX_SVC[Outbox Service]
        USER_SVC[User Service]
    end
    
    subgraph "Data Access Layer (Repository)"
        AUTH_REPO[Auth Repository]
        SESS_REPO[Session SQL Repository]
        DOMAIN_REPO[Domain SQL Repository]
        OUTBOX_REPO[Outbox SQL Repository]
    end
    
    subgraph "Database (PostgreSQL)"
        DB[(PostgreSQL)]
        US[Users Table]
        SS[Sessions Table]
        DE[Domain Events Table]
        TO[Transactional Outbox]
        AE[Audit Events Table]
    end
    
    %% Connections
    WEB --> API
    MW_AUTH --> AUTH_SVC
    MW_AUTH --> SESS_SVC
    
    AUTH_SVC --> USER_SVC
    AUTH_SVC --> SESS_SVC
    
    AUTH_SVC --> AUTH_REPO
    SESS_SVC --> SESS_REPO
    DOMAIN_SVC --> DOMAIN_REPO
    OUTBOX_SVC --> OUTBOX_REPO
    
    AUTH_REPO --> DB
    SESS_REPO --> DB
    DOMAIN_REPO --> DB
    OUTBOX_REPO --> DB
```

---

## 2. API Call Flow (Middleware Chain & Domain Mutating Actions)

All authenticated and state-mutating requests go through a strict middleware chain that enforces observability, audit logging, and authorization, followed by an ACID database transaction.

```mermaid
sequenceDiagram
    participant Client
    participant Router as HTTP Router
    participant CORS as CORS Middleware
    participant Audit as Audit Middleware
    participant Auth as Auth Middleware
    participant Handler
    participant Service
    participant DB as PostgreSQL

    Client->>Router: HTTP Request (e.g. POST /api/action)
    Router->>CORS: Preflight / Request Validation
    CORS-->>Router: Allow Origin
    Router->>Audit: Record Request Start (Assign Correlation ID)
    Router->>Auth: Validate Session Cookie
    Auth->>DB: Query Session Active Status
    DB-->>Auth: Session Valid
    Auth-->>Router: Inject UserContext
    Router->>Handler: ServeHTTP
    Handler->>Service: Execute Domain Logic
    
    %% Transactional Block
    Service->>DB: Begin ACID Tx
    Service->>DB: Mutate Domain State
    Service->>DB: Insert Domain Event
    Service->>DB: Insert Transactional Outbox Event
    Service->>DB: Commit Tx
    
    %% Audit Log Completion
    Service-->>Handler: Return Result
    Handler-->>Router: HTTP Response (200 OK)
    Router->>Audit: Record Request Finish (Async Audit Event Logged)
    Audit->>DB: Insert Audit Event
    Router-->>Client: HTTP Response
```

---

## 3. User Flows

### A. User Registration Flow

```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant AuthSvc as Auth Service
    participant DB as PostgreSQL

    User->>Web: Submit Registration Form (Name, DOB, Email, Mobile)
    Web->>API: POST /register
    API->>AuthSvc: RegisterUser(request)
    AuthSvc->>AuthSvc: Validate Input & Ensure uniqueness
    AuthSvc->>AuthSvc: Hash Password (Argon2id)
    
    AuthSvc->>DB: Begin Tx
    AuthSvc->>DB: INSERT INTO users
    AuthSvc->>DB: INSERT INTO domain_events (UserRegistered)
    AuthSvc->>DB: INSERT INTO transactional_outbox
    AuthSvc->>DB: Commit Tx
    
    AuthSvc-->>API: Registration Success
    API-->>Web: Return 201 Created
    Web-->>User: Prompt to Login
```

### B. User Login & Session Creation

```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant AuthSvc as Auth Service
    participant SessSvc as Session Service
    participant DB as PostgreSQL

    User->>Web: Submit Credentials (Username/Password)
    Web->>API: POST /login
    API->>AuthSvc: Authenticate(credentials)
    
    AuthSvc->>DB: Query User by Username
    DB-->>AuthSvc: Return User Record + Password Hash
    AuthSvc->>AuthSvc: Verify Argon2id Hash
    
    AuthSvc->>SessSvc: CreateSession(userID, IP, UserAgent)
    SessSvc->>DB: Begin Tx
    SessSvc->>DB: INSERT INTO sessions (status = active)
    SessSvc->>DB: Commit Tx
    SessSvc-->>AuthSvc: Session Created
    
    AuthSvc-->>API: Login Success (Token/Session ID)
    API-->>Web: Set-Cookie: session_id (HttpOnly, Secure)
    Web-->>User: Redirect to Authenticated Dashboard
```

### C. Session Revocation (Logout)

```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant SessSvc as Session Service
    participant DB as PostgreSQL

    User->>Web: Click Logout
    Web->>API: POST /logout (Includes Session Cookie)
    API->>SessSvc: RevokeSession(sessionID)
    
    SessSvc->>DB: UPDATE sessions SET status = 'revoked', revoked_at = NOW()
    DB-->>SessSvc: Success
    
    SessSvc-->>API: Revocation Success
    API-->>Web: Clear-Cookie & 200 OK
    Web-->>User: Redirect to Public Home/Login
```

---

## 4. Key Architectural Patterns Employed

1. **Transactional Outbox**: All domain events (e.g., `UserRegistered`, `PasswordChanged`) are saved to both a `domain_events` log and a `transactional_outbox` table within the **same database transaction** as the entity mutation. This guarantees that events are never lost and can be processed asynchronously by a background worker for notifications, webhooks, or syncing to external systems.
2. **Session-based Authentication**: Instead of stateless JWTs, sessions are stateful and persisted in the `sessions` table. This provides immediate revocation capabilities, concurrent session limits, and tight IP/Device tracking for security audits.
3. **Structured Audit Logging**: Operations produce immutable audit records linking `request_id`, `user_id`, `session_id`, `action`, and `result` (Success/Failure) to ensure 100% traceability for financial data security.
