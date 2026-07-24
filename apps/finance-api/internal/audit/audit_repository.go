package audit

import (
	"context"
	"database/sql"
	"errors"
)

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

// InsertTx inserts an audit event using the provided transaction.
func (r *SQLRepository) InsertTx(ctx context.Context, tx *sql.Tx, e AuditEvent) error {
	q := `INSERT INTO audit_events (id, request_id, user_id, session_id, resource, resource_id, action, result, timestamp, user_agent, client_ip) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := tx.ExecContext(ctx, q, e.ID, e.RequestID, nullableString(e.UserID), nullableString(e.SessionID), e.Resource, nullableString(e.ResourceID), e.Action, e.Result, e.Timestamp, e.UserAgent, nullableString(e.ClientIP))
	return err
}

func (r *SQLRepository) Update(ctx context.Context, e AuditEvent) error {
	return errors.New("audit records are append-only")
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	return errors.New("audit records are append-only")
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
