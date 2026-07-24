package audit

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuditRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)

	evt := AuditEvent{
		ID:        "evt-1",
		RequestID: "req-1",
		UserID:    "user-1",
		SessionID: "sess-1",
		Resource:  "user",
		Action:    "create",
		Result:    "success",
		Timestamp: time.Now(),
		UserAgent: "test",
		ClientIP:  "127.0.0.1",
	}

	mock.ExpectExec("INSERT INTO audit_events").WithArgs(
		evt.ID, evt.RequestID, evt.UserID, evt.SessionID, evt.Resource, nil, evt.Action, evt.Result, evt.Timestamp, evt.UserAgent, evt.ClientIP,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(context.Background(), evt)
	if err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}

	mock.ExpectBegin()
	tx, _ := db.Begin()
	mock.ExpectExec("INSERT INTO audit_events").WithArgs(
		evt.ID, evt.RequestID, evt.UserID, evt.SessionID, evt.Resource, nil, evt.Action, evt.Result, evt.Timestamp, evt.UserAgent, evt.ClientIP,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.InsertTx(context.Background(), tx, evt)
	if err != nil {
		t.Errorf("error was not expected while inserting tx: %s", err)
	}
	tx.Commit()

	if err := repo.Update(context.Background(), evt); err == nil {
		t.Error("expected error on update")
	}

	if err := repo.Delete(context.Background(), "id"); err == nil {
		t.Error("expected error on delete")
	}

	// test null string
	if nullableString("") != nil {
		t.Error("expected nil")
	}
}
