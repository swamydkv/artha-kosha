package outbox
package outbox

import (
    "context"
    "database/sql"
    "time"
)

type OutboxEntry struct {
    ID            string
    DomainEventID string
    EventType     string
    Payload       []byte
    CreatedAt     time.Time
    ProcessingStatus string
}

type Repository interface {
    Insert(ctx context.Context, e OutboxEntry) error
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Enqueue(ctx context.Context, e OutboxEntry) error {
    if e.CreatedAt.IsZero() {
        e.CreatedAt = time.Now().UTC()
    }
    return s.repo.Insert(ctx, e)
}

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) Insert(ctx context.Context, e OutboxEntry) error {
    q := `INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count) VALUES ($1,$2,$3,$4,$5,$6,$7)`
    _, err := r.db.ExecContext(ctx, q, e.ID, e.DomainEventID, e.EventType, e.Payload, e.CreatedAt, e.ProcessingStatus, 0)
    return err
}
