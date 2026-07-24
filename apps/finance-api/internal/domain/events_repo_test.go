package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	e := DomainEvent{
		ID:               "1",
		EventType:        "test",
		AggregateID:      "2",
		AggregateType:    "agg",
		EventData:        []byte("{}"),
		Timestamp:        time.Now(),
		ProcessingStatus: "pending",
		RetryCount:       0,
	}

	mock.ExpectExec("INSERT INTO domain_events").
		WithArgs(e.ID, e.EventType, e.AggregateID, e.AggregateType, e.EventData, e.Timestamp, e.ProcessingStatus, e.RetryCount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(context.Background(), e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_InsertTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	e := DomainEvent{
		ID:               "1",
		EventType:        "test",
		AggregateID:      "2",
		AggregateType:    "agg",
		EventData:        []byte("{}"),
		Timestamp:        time.Now(),
		ProcessingStatus: "pending",
		RetryCount:       0,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain_events").
		WithArgs(e.ID, e.EventType, e.AggregateID, e.AggregateType, e.EventData, e.Timestamp, e.ProcessingStatus, e.RetryCount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	err = repo.InsertTx(context.Background(), tx, e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_GetPending_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "event_type", "aggregate_id", "aggregate_type", "event_data", "timestamp", "processing_status", "retry_count"}).
		AddRow("1", "test", "2", "agg", []byte("{}"), now, "pending", 0)

	mock.ExpectQuery("SELECT id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count FROM domain_events").
		WithArgs(10).
		WillReturnRows(rows)

	events, err := repo.GetPending(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ID != "1" {
		t.Errorf("expected id 1, got %s", events[0].ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_GetPending_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectQuery("SELECT id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count FROM domain_events").
		WithArgs(10).
		WillReturnError(errors.New("db error"))

	_, err = repo.GetPending(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSQLRepository_GetPending_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	// Return a row with missing column to trigger scan error
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("1")

	mock.ExpectQuery("SELECT id, event_type, aggregate_id, aggregate_type, event_data, timestamp, processing_status, retry_count FROM domain_events").
		WithArgs(10).
		WillReturnRows(rows)

	_, err = repo.GetPending(context.Background(), 10)
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
