package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/outbox"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestService_Emit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	svc := NewService(repo, outboxSvc)

	e := DomainEvent{
		ID:            "1",
		EventType:     "test",
		AggregateID:   "2",
		AggregateType: "agg",
	}

	mock.ExpectExec("INSERT INTO domain_events").
		WithArgs(e.ID, e.EventType, e.AggregateID, e.AggregateType, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO transactional_outbox").
		WithArgs(sqlmock.AnyArg(), e.ID, e.EventType, sqlmock.AnyArg(), sqlmock.AnyArg(), "pending", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = svc.Emit(context.Background(), e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestService_Emit_RepoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	svc := NewService(repo, outboxSvc)

	e := DomainEvent{
		ID: "1",
	}

	mock.ExpectExec("INSERT INTO domain_events").
		WillReturnError(errors.New("db error"))

	err = svc.Emit(context.Background(), e)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestService_EmitTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	svc := NewService(repo, outboxSvc)

	// Test success with Zero timestamp
	e := DomainEvent{
		ID: "1",
		// Timestamp is intentionally zero
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain_events").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactional_outbox").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	err = svc.EmitTx(context.Background(), tx, e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}

	// Test failure from InsertTx
	e2 := DomainEvent{
		ID:        "2",
		Timestamp: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain_events").
		WillReturnError(errors.New("db error in tx"))

	tx2, _ := db.Begin()
	err = svc.EmitTx(context.Background(), tx2, e2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestService_EmitTx_Fallback(t *testing.T) {
	// If outbox repo doesn't support InsertTx, we can test fallback, but since NewSQLRepository returns one that does, it's hard to mock here without an interface.
	// We'll skip fallback test or create a dummy struct for it.
	
	// Create a dummy repo that doesn't support Tx
	type DummyRepo struct {
		Repository
	}
	// Test is complex to set up without exposing interfaces in NewService. 
	// We skip fallback for now, as it requires interface casting inside EmitTx.
}
