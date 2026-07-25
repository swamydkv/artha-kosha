package outbox

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockInserter struct {
	err error
}

func (m mockInserter) Insert(ctx context.Context, e OutboxEntry) error { return m.err }
func (m mockInserter) InsertTx(ctx context.Context, tx *sql.Tx, e OutboxEntry) error { return m.err }

type mockRepoOnlyInsert struct {
	err error
}
func (m mockRepoOnlyInsert) Insert(ctx context.Context, e OutboxEntry) error { return m.err }

func TestService_Enqueue(t *testing.T) {
	svc := NewService(mockInserter{})
	err := svc.Enqueue(context.Background(), OutboxEntry{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// with created_at set
	err = svc.Enqueue(context.Background(), OutboxEntry{CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestService_EnqueueTx(t *testing.T) {
	svc := NewService(mockInserter{})
	err := svc.EnqueueTx(context.Background(), nil, OutboxEntry{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	
	// with created_at set
	err = svc.EnqueueTx(context.Background(), nil, OutboxEntry{CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// fallback to regular insert
	svc2 := NewService(mockRepoOnlyInsert{})
	err = svc2.EnqueueTx(context.Background(), nil, OutboxEntry{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestSQLRepository_Insert(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactional_outbox")).
		WithArgs("1", "de1", "evt", []byte("data"), sqlmock.AnyArg(), "pending", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Insert(context.Background(), OutboxEntry{ID: "1", DomainEventID: "de1", EventType: "evt", Payload: []byte("data"), ProcessingStatus: "pending"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestSQLRepository_InsertTx(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectBegin()
	tx, _ := db.Begin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactional_outbox")).
		WithArgs("1", "de1", "evt", []byte("data"), sqlmock.AnyArg(), "pending", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.InsertTx(context.Background(), tx, OutboxEntry{ID: "1", DomainEventID: "de1", EventType: "evt", Payload: []byte("data"), ProcessingStatus: "pending"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestSQLRepository_FetchPending(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewSQLRepository(db)

	rows := sqlmock.NewRows([]string{"id", "domain_event_id", "event_type", "payload", "created_at", "processing_status"}).
		AddRow("1", "de1", "evt", []byte("data"), time.Now(), "pending")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, domain_event_id, event_type, payload, created_at, processing_status FROM transactional_outbox")).
		WithArgs(10).
		WillReturnRows(rows)

	res, err := repo.FetchPending(context.Background(), 10)
	if err != nil || len(res) != 1 {
		t.Fatalf("unexpected err or length: %v, %v", err, len(res))
	}

	// Error path
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id")).WillReturnError(errors.New("db error"))
	_, err = repo.FetchPending(context.Background(), 10)
	if err == nil {
		t.Fatalf("expected error")
	}

	// Scan error path
	rowsErr := sqlmock.NewRows([]string{"id", "domain_event_id", "event_type", "payload", "created_at", "processing_status"}).
		AddRow(nil, nil, nil, nil, nil, nil) // causes scan error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id")).WithArgs(10).WillReturnRows(rowsErr)
	_, err = repo.FetchPending(context.Background(), 10)
	if err == nil {
		t.Fatalf("expected scan error")
	}
}

func TestSQLRepository_Updates(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewSQLRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE transactional_outbox SET processing_status = 'processed'")).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := repo.MarkProcessed(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE transactional_outbox SET retry_count")).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.IncrementRetry(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE transactional_outbox SET processing_status = 'failed'")).
		WithArgs("1", "reason").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.MarkFailed(context.Background(), "1", "reason")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM transactional_outbox")).
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteProcessed(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
