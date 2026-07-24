package domain

import (
	"context"
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

func generateID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf))
}
