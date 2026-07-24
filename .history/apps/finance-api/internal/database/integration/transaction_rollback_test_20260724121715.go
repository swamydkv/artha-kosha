package integration

import (
    "context"
    "errors"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"

    "artha-kosha/apps/finance-api/internal/database"
)

func TestWithTx_RollsBackOnError(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create sqlmock: %v", err)
    }
    defer db.Close()

    // Expect a transaction begin
    mock.ExpectBegin()

    // Expect an Exec inside the transaction
    mock.ExpectExec("INSERT INTO test_table").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

    // Expect rollback because our fn will return an error
    mock.ExpectRollback()

    err = database.WithTx(context.Background(), db, func(tx *sql.Tx) error {
        if _, e := tx.ExecContext(context.Background(), "INSERT INTO test_table (id) VALUES (?)", 1); e != nil {
            return e
        }
        return errors.New("force rollback")
    })
    if err == nil {
        t.Fatalf("expected error from WithTx, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet mock expectations: %v", err)
    }
}
