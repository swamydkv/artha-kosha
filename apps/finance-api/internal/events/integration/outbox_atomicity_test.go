package integration

import (
	"context"
	"errors"
	"testing"

	"artha-kosha/apps/finance-api/internal/domain"
	"artha-kosha/apps/finance-api/internal/outbox"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOutboxAtomicity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	// We start a mock transaction
	mock.ExpectBegin()
	// Domain event insert succeeds (8 args)
	mock.ExpectExec("INSERT INTO domain_events").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	// Outbox insert fails (7 args)
	mock.ExpectExec("INSERT INTO transactional_outbox").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnError(errors.New("db error"))

	repo := domain.NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	svc := domain.NewService(repo, outboxSvc)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	err = svc.EmitTx(context.Background(), tx, domain.DomainEvent{
		ID:            "evt-123",
		EventType:     "TEST",
		AggregateID:   "agg-123",
		AggregateType: "test",
	})

	if err == nil {
		t.Error("Expected error when outbox insertion fails, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}
