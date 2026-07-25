package domain

import (
	"context"
	"database/sql"
	"time"

	"artha-kosha/apps/finance-api/internal/outbox"
	"github.com/google/uuid"
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
		ID:               uuid.New().String(),
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
	if err := s.repo.InsertTx(ctx, tx, e); err != nil {
		return err
	}

	// prepare outbox entry
	oe := outbox.OutboxEntry{
		ID:               uuid.New().String(),
		DomainEventID:    e.ID,
		EventType:        e.EventType,
		Payload:          e.EventData,
		CreatedAt:        time.Now().UTC(),
		ProcessingStatus: "pending",
	}

	// use transactional enqueue
	return s.outbox.EnqueueTx(ctx, tx, oe)
}

