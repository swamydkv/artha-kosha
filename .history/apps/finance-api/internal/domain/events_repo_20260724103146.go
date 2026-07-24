package domain

import (
	"context"
	"database/sql"
	"time"
)

type DomainEvent struct {
	ID               string
	EventType        string
	AggregateID      string
	AggregateType    string
	EventData        []byte
	Timestamp        time.Time
	ProcessingStatus string
	RetryCount       int
}

type Repository interface {
	Insert(ctx context.Context, e DomainEvent) error
	GetPending(ctx context.Context, limit int) ([]DomainEvent, error)
}

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) Insert(ctx context.Context, e DomainEvent) error {
	q := `INSERT INTO domain_events (id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.ExecContext(ctx, q, e.ID, e.EventType, e.AggregateID, e.AggregateType, e.EventData, e.Timestamp, e.ProcessingStatus, e.RetryCount)
	return err
}

func (r *SQLRepository) GetPending(ctx context.Context, limit int) ([]DomainEvent, error) {
	q := `SELECT id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count FROM domain_events WHERE processing_status = 'pending' ORDER BY timestamp ASC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DomainEvent
	for rows.Next() {
		var e DomainEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.AggregateID, &e.AggregateType, &e.EventData, &e.Timestamp, &e.ProcessingStatus, &e.RetryCount); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
