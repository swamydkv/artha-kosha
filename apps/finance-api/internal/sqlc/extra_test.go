package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestQueries_RowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	q := New(db)
	ctx := context.Background()
	expectedErr := sqlmock.ErrCancelled

	// GetAuditEventsByRequestID
	rows1 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).
		AddRow(nil, nil, nil, nil, nil, nil, nil, nil, time.Now(), nil, nil).
		RowError(0, expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows1)
	q.GetAuditEventsByRequestID(ctx, "")

	// FetchPendingDomainEvents
	rows2 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9"}).
		AddRow(nil, nil, nil, nil, nil, time.Now(), nil, 0, nil, nil).
		RowError(0, expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows2)
	q.FetchPendingDomainEvents(ctx, 0)

	// FetchPendingOutbox
	rows3 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).
		AddRow(nil, nil, nil, nil, time.Now(), nil, 0, nil, nil, nil, nil).
		RowError(0, expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows3)
	q.FetchPendingOutbox(ctx, 0)
}

func TestQueries_RowsCloseErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	q := New(db)
	ctx := context.Background()
	expectedErr := sqlmock.ErrCancelled

	rows1 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).
		AddRow("00000000-0000-0000-0000-000000000000", "req", nil, nil, "res", nil, "act", "success", time.Now(), nil, nil).
		CloseError(expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows1)
	q.GetAuditEventsByRequestID(ctx, "")

	rows2 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).
		AddRow("00000000-0000-0000-0000-000000000000", "USER_REGISTERED", "agg", "aggType", []byte("{}"), time.Now(), "pending", 0).
		CloseError(expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows2)
	q.FetchPendingDomainEvents(ctx, 0)

	rows3 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6"}).
		AddRow("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "evt", []byte("{}"), time.Now(), "pending", 0).
		CloseError(expectedErr)
	mock.ExpectQuery(".*").WillReturnRows(rows3)
	q.FetchPendingOutbox(ctx, 0)
}

func TestQueries_RowsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	q := New(db)
	ctx := context.Background()

	rows1 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).
		AddRow("00000000-0000-0000-0000-000000000000", "req", nil, nil, "res", nil, "act", "success", time.Now(), nil, nil)
	mock.ExpectQuery(".*").WillReturnRows(rows1)
	q.GetAuditEventsByRequestID(ctx, "")

	rows2 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).
		AddRow("00000000-0000-0000-0000-000000000000", "USER_REGISTERED", "agg", "aggType", []byte("{}"), time.Now(), "pending", 0)
	mock.ExpectQuery(".*").WillReturnRows(rows2)
	q.FetchPendingDomainEvents(ctx, 0)

	rows3 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6"}).
		AddRow("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "evt", []byte("{}"), time.Now(), "pending", 0)
	mock.ExpectQuery(".*").WillReturnRows(rows3)
	q.FetchPendingOutbox(ctx, 0)
}

func TestQueries_RowsScanErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	q := New(db)
	ctx := context.Background()

	rows1 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}).
		AddRow(5, "req", nil, nil, "res", nil, "act", "success", time.Now(), nil, nil)
	mock.ExpectQuery(".*").WillReturnRows(rows1)
	q.GetAuditEventsByRequestID(ctx, "")

	rows2 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6", "col7"}).
		AddRow(5, "USER_REGISTERED", "agg", "aggType", []byte("{}"), time.Now(), "pending", 0)
	mock.ExpectQuery(".*").WillReturnRows(rows2)
	q.FetchPendingDomainEvents(ctx, 0)

	rows3 := sqlmock.NewRows([]string{"col0", "col1", "col2", "col3", "col4", "col5", "col6"}).
		AddRow(5, "00000000-0000-0000-0000-000000000000", "evt", []byte("{}"), time.Now(), "pending", 0)
	mock.ExpectQuery(".*").WillReturnRows(rows3)
	q.FetchPendingOutbox(ctx, 0)
}
