package outbox

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockRepo struct {
	insertErr error
	insertTxErr error
}

func (m *mockRepo) Insert(ctx context.Context, e OutboxEntry) error {
	return m.insertErr
}

func (m *mockRepo) InsertTx(ctx context.Context, tx *sql.Tx, e OutboxEntry) error {
	return m.insertTxErr
}

func TestService_Enqueue(t *testing.T) {
	svc := NewService(&mockRepo{})
	err := svc.Enqueue(context.Background(), OutboxEntry{})
	if err != nil {
		t.Error("expected no error")
	}

	err = svc.Enqueue(context.Background(), OutboxEntry{CreatedAt: time.Now()})
	if err != nil {
		t.Error("expected no error")
	}
}

func TestService_EnqueueTx(t *testing.T) {
	svc := NewService(&mockRepo{})
	err := svc.EnqueueTx(context.Background(), nil, OutboxEntry{})
	if err != nil {
		t.Error("expected no error")
	}

	err = svc.EnqueueTx(context.Background(), nil, OutboxEntry{CreatedAt: time.Now()})
	if err != nil {
		t.Error("expected no error")
	}
	
	// Test fallback to Insert when InsertTx is not implemented
	type mockRepoNoTx struct {
		insertErr error
	}
	// can't easily dynamically cast without full type, let's just make a new struct
}
// define fallback
type fallbackRepoType struct{}
func (f *fallbackRepoType) Insert(ctx context.Context, e OutboxEntry) error { return nil }

func TestService_EnqueueTxFallback(t *testing.T) {
	svc := NewService(&fallbackRepoType{})
	err := svc.EnqueueTx(context.Background(), nil, OutboxEntry{})
	if err != nil {
		t.Error("expected no error")
	}
}

func TestSQLRepository(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewSQLRepository(db)

	// Insert
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.Insert(context.Background(), OutboxEntry{})
	if err != nil {
		t.Error("expected no error")
	}

	// InsertTx
	mock.ExpectBegin()
	tx, _ := db.Begin()
	mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.InsertTx(context.Background(), tx, OutboxEntry{})
	if err != nil {
		t.Error("expected no error")
	}

	// FetchPending
	rows := sqlmock.NewRows([]string{"id", "domain_event_id", "event_type", "payload", "created_at", "processing_status"}).
		AddRow("id1", "de1", "t", []byte{}, time.Now(), "pending")
	mock.ExpectQuery("SELECT id").WillReturnRows(rows)
	entries, err := repo.FetchPending(context.Background(), 10)
	if err != nil || len(entries) != 1 {
		t.Error("expected 1 entry")
	}
	
	mock.ExpectQuery("SELECT id").WillReturnError(sql.ErrNoRows)
	_, err = repo.FetchPending(context.Background(), 10)
	if err == nil {
		t.Error("expected error")
	}

	rowsErr := sqlmock.NewRows([]string{"id", "domain_event_id", "event_type", "payload", "created_at", "processing_status"}).
		AddRow("id1", "de1", "t", []byte{}, time.Now(), "pending").RowError(0, sql.ErrNoRows)
	mock.ExpectQuery("SELECT id").WillReturnRows(rowsErr)
	_, err = repo.FetchPending(context.Background(), 10)
	if err == nil {
		t.Error("expected error")
	}

	// MarkProcessed
	mock.ExpectExec("UPDATE transactional_outbox SET processing_status = 'processed'").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.MarkProcessed(context.Background(), "id1")
	if err != nil {
		t.Error("expected no error")
	}

	// IncrementRetry
	mock.ExpectExec("UPDATE transactional_outbox SET retry_count").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.IncrementRetry(context.Background(), "id1")
	if err != nil {
		t.Error("expected no error")
	}

	// MarkFailed
	mock.ExpectExec("UPDATE transactional_outbox SET processing_status = 'failed'").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.MarkFailed(context.Background(), "id1", "reason")
	if err != nil {
		t.Error("expected no error")
	}
}
