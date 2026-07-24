package domain

import (
	"context"
	"database/sql"
	"time"

	"artha-kosha/apps/finance-api/internal/outbox"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Service struct {
	repo   *SQLRepository
	outbox *outbox.Service
}

func NewService(repo *SQLRepository, out *outbox.Service) *Service {
	return &Service{repo: repo, outbox: out}
}

// Emit records a domain event and enqueues an outbox entry for delivery.
func (s *Service) Emit(ctx context.Context, e DomainEvent) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	if err := s.repo.Insert(ctx, e); err != nil {
		return err
	}
	// enqueue outbox entry linking to domain event
	oe := outbox.OutboxEntry{
		ID:               generateID("outbox"),
		DomainEventID:    e.ID,
		EventType:        e.EventType,
		Payload:          e.EventData,
		CreatedAt:        time.Now().UTC(),
		ProcessingStatus: "pending",
	}
	return s.outbox.Enqueue(ctx, oe)
}

// EmitTx records the domain event and enqueues an outbox entry within the provided transaction.
// It will use repository transactional insert methods when available.
func (s *Service) EmitTx(ctx context.Context, tx *sql.Tx, e DomainEvent) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	// attempt to use repository transaction insert
	if inserter, ok := interface{}(s.repo).(interface {
		InsertTx(context.Context, *sql.Tx, DomainEvent) error
	}); ok {
		if err := inserter.InsertTx(ctx, tx, e); err != nil {
			return err
		}
	} else {
		// fallback to non-transactional insert
		if err := s.repo.Insert(ctx, e); err != nil {
			return err
		}
	}

	// prepare outbox entry
	oe := outbox.OutboxEntry{
		ID:               generateID("outbox"),
		DomainEventID:    e.ID,
		EventType:        e.EventType,
		Payload:          e.EventData,
		CreatedAt:        time.Now().UTC(),
		ProcessingStatus: "pending",
	}

	// best-effort: if outbox service's repository supports InsertTx, use that via an exported method on outbox package
	if inserter, ok := interface{}(s.outbox).(interface {
		EnqueueTx(context.Context, *sql.Tx, outbox.OutboxEntry) error
	}); ok {
		return inserter.EnqueueTx(ctx, tx, oe)
	}

	// fallback: use non-transactional enqueue (best-effort)
	return s.outbox.Enqueue(ctx, oe)
}

func generateID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf))
}
