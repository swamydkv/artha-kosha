package audit

import (
    "context"
    "database/sql"
    "time"
)

type AuditEvent struct {
    ID        string
    RequestID string
    UserID    string
    SessionID string
    Resource  string
    ResourceID string
    Action    string
    Result    string
    Timestamp time.Time
    UserAgent string
    ClientIP  string
}

type Repository interface {
    Insert(ctx context.Context, e AuditEvent) error
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, e AuditEvent) error {
    if e.Timestamp.IsZero() {
        e.Timestamp = time.Now().UTC()
    }
    return s.repo.Insert(ctx, e)
}

// SQLRepository is a simple repo backed by *sql.DB using the audit_events queries.
type SQLRepository struct {
    db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) Insert(ctx context.Context, e AuditEvent) error {
    q := `INSERT INTO audit_events (id, request_id, user_id, session_id, resource, resource_id, action, result, timestamp, user_agent, client_ip) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
    _, err := r.db.ExecContext(ctx, q, e.ID, e.RequestID, nullableString(e.UserID), nullableString(e.SessionID), e.Resource, nullableString(e.ResourceID), e.Action, e.Result, e.Timestamp, e.UserAgent, nullableString(e.ClientIP))
    return err
}

func nullableString(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}
