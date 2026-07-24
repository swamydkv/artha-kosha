package audit

import (
	"context"
	"database/sql"
	"time"
)

type AuditEvent struct {
	ID         string
	RequestID  string
	UserID     string
	SessionID  string
	Resource   string
	ResourceID string
	Action     string
	Result     string
	Timestamp  time.Time
	UserAgent  string
	ClientIP   string
}

type Repository interface {
	Insert(ctx context.Context, e AuditEvent) error
}

// TxInserter is implemented by repositories that support transactional inserts.
type TxInserter interface {
	InsertTx(ctx context.Context, tx *sql.Tx, e AuditEvent) error
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

// RecordTx attempts to insert the audit event within the provided transaction.
// If the underlying repository implements InsertTx, it will be used. Otherwise
// falls back to non-transactional Record (best-effort).
func (s *Service) RecordTx(ctx context.Context, tx *sql.Tx, e AuditEvent) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	if inserter, ok := s.repo.(TxInserter); ok {
		return inserter.InsertTx(ctx, tx, e)
	}
	// fallback to non-transactional insert
	return s.repo.Insert(ctx, e)
}

