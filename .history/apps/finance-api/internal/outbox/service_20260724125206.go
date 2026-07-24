package outbox

import (
	"context"
	"database/sql"
	"time"
)

type OutboxEntry struct {
	ID               string
	DomainEventID    string
	EventType        string
	Payload          []byte
	CreatedAt        time.Time
	ProcessingStatus string
}

type Repository interface {
	Insert(ctx context.Context, e OutboxEntry) error
}

// ReadWriter defines operations the worker needs to fetch and update outbox entries.
type ReadWriter interface {
	Insert(ctx context.Context, e OutboxEntry) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEntry, error)
	MarkProcessed(ctx context.Context, id string) error
	IncrementRetry(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, reason string) error
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Enqueue(ctx context.Context, e OutboxEntry) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	return s.repo.Insert(ctx, e)
}

// EnqueueTx attempts to insert the outbox entry within the provided transaction.
// If the underlying repository implements InsertTx, it will be used. Otherwise
// falls back to non-transactional Insert (best-effort).
func (s *Service) EnqueueTx(ctx context.Context, tx *sql.Tx, e OutboxEntry) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	if inserter, ok := s.repo.(interface{ InsertTx(context.Context, *sql.Tx, OutboxEntry) error }); ok {
		return inserter.InsertTx(ctx, tx, e)
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

// InsertTx inserts an outbox entry using the provided transaction.
func (r *SQLRepository) InsertTx(ctx context.Context, tx *sql.Tx, e OutboxEntry) error {
	q := `INSERT INTO transactional_outbox (id, domain_event_id, event_type, payload, created_at, processing_status, retry_count) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := tx.ExecContext(ctx, q, e.ID, e.DomainEventID, e.EventType, e.Payload, e.CreatedAt, e.ProcessingStatus, 0)
	return err
}

func (r *SQLRepository) FetchPending(ctx context.Context, limit int) ([]OutboxEntry, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, domain_event_id, event_type, payload, created_at, processing_status FROM transactional_outbox WHERE processing_status = 'pending' ORDER BY created_at ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []OutboxEntry
	for rows.Next() {
		var e OutboxEntry
		if err := rows.Scan(&e.ID, &e.DomainEventID, &e.EventType, &e.Payload, &e.CreatedAt, &e.ProcessingStatus); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, rows.Err()
}

func (r *SQLRepository) MarkProcessed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE transactional_outbox SET processing_status = 'processed', processed_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *SQLRepository) IncrementRetry(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE transactional_outbox SET retry_count = COALESCE(retry_count,0)+1, last_retry_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *SQLRepository) MarkFailed(ctx context.Context, id string, reason string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE transactional_outbox SET processing_status = 'failed', last_error = $2, processed_at = NOW() WHERE id = $1`, id, reason)
	return err
}
