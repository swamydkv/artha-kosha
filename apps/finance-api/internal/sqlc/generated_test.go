package sqlc

import (
	"context"
	"testing"


	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestQueries_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	q := New(db)
	ctx := context.Background()

	// Test GetAuditEventsByRequestID
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetAuditEventsByRequestID(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetAuditEventsByRequestID(ctx, "")

	// Test InsertAuditEvent
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertAuditEvent(ctx, InsertAuditEventParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertAuditEvent(ctx, InsertAuditEventParams{})

	// Test InsertBudget
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertBudget(ctx, InsertBudgetParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertBudget(ctx, InsertBudgetParams{})

	// Test FetchPendingDomainEvents
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
	q.FetchPendingDomainEvents(ctx, 0)
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.FetchPendingDomainEvents(ctx, 0)

	// Test InsertDomainEvent
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertDomainEvent(ctx, InsertDomainEventParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertDomainEvent(ctx, InsertDomainEventParams{})

	// Test MarkDomainEventProcessed
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.MarkDomainEventProcessed(ctx, uuid.Nil)
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.MarkDomainEventProcessed(ctx, uuid.Nil)

	// Test FetchPendingOutbox
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6"}).AddRow(nil, nil, nil, nil, nil, nil, nil))
	q.FetchPendingOutbox(ctx, 0)
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.FetchPendingOutbox(ctx, 0)

	// Test IncrementOutboxRetry
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.IncrementOutboxRetry(ctx, uuid.Nil)
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.IncrementOutboxRetry(ctx, uuid.Nil)

	// Test InsertOutboxEntry
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertOutboxEntry(ctx, InsertOutboxEntryParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertOutboxEntry(ctx, InsertOutboxEntryParams{})

	// Test MarkOutboxFailed
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.MarkOutboxFailed(ctx, MarkOutboxFailedParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.MarkOutboxFailed(ctx, MarkOutboxFailedParams{})

	// Test MarkOutboxProcessed
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.MarkOutboxProcessed(ctx, uuid.Nil)
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.MarkOutboxProcessed(ctx, uuid.Nil)

	// Test CreateSession
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.CreateSession(ctx, CreateSessionParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.CreateSession(ctx, CreateSessionParams{})

	// Test GetSessionByID
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetSessionByID(ctx, uuid.Nil)
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetSessionByID(ctx, uuid.Nil)

	// Test RevokeAllSessionsByUser
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.RevokeAllSessionsByUser(ctx, RevokeAllSessionsByUserParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.RevokeAllSessionsByUser(ctx, RevokeAllSessionsByUserParams{})

	// Test RevokeSession
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.RevokeSession(ctx, RevokeSessionParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.RevokeSession(ctx, RevokeSessionParams{})

	// Test UpdateSessionActivity
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.UpdateSessionActivity(ctx, UpdateSessionActivityParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.UpdateSessionActivity(ctx, UpdateSessionActivityParams{})

	// Test InsertTransaction
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertTransaction(ctx, InsertTransactionParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertTransaction(ctx, InsertTransactionParams{})

	// Test CheckEmailExists
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0"}).AddRow(nil))
	q.CheckEmailExists(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.CheckEmailExists(ctx, "")

	// Test CheckMobileNumberExists
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0"}).AddRow(nil))
	q.CheckMobileNumberExists(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.CheckMobileNumberExists(ctx, "")

	// Test CheckUsernameExists
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0"}).AddRow(nil))
	q.CheckUsernameExists(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.CheckUsernameExists(ctx, "")

	// Test CreateUser
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4"}).AddRow(nil, nil, nil, nil, nil))
	q.CreateUser(ctx, CreateUserParams{})
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.CreateUser(ctx, CreateUserParams{})

	// Test GetUserByEmail
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetUserByEmail(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetUserByEmail(ctx, "")

	// Test GetUserByID
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetUserByID(ctx, uuid.Nil)
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetUserByID(ctx, uuid.Nil)

	// Test GetUserByMobileNumber
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetUserByMobileNumber(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetUserByMobileNumber(ctx, "")

	// Test GetUserByUsername
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))
	q.GetUserByUsername(ctx, "")
	mock.ExpectQuery(".*").WillReturnError(sqlmock.ErrCancelled)
	q.GetUserByUsername(ctx, "")

	// Test InsertAccount
	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(1, 1))
	q.InsertAccount(ctx, InsertAccountParams{})
	mock.ExpectExec(".*").WillReturnError(sqlmock.ErrCancelled)
	q.InsertAccount(ctx, InsertAccountParams{})
}
