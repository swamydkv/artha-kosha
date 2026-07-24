package integration

import (
    "context"
    "errors"
    "testing"
    "database/sql"

    "github.com/DATA-DOG/go-sqlmock"

    "artha-kosha/apps/finance-api/internal/database"
)

func TestWithTx_AtomicMultiTableOperations(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create sqlmock: %v", err)
    }
    defer db.Close()

    t.Run("commit on success", func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO table_a").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectExec("INSERT INTO table_b").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectCommit()

        err := database.WithTx(context.Background(), db, func(tx *sql.Tx) error {
            if _, e := tx.ExecContext(context.Background(), "INSERT INTO table_a (id) VALUES (?)", 1); e != nil {
                return e
            }
            if _, e := tx.ExecContext(context.Background(), "INSERT INTO table_b (id) VALUES (?)", 2); e != nil {
                return e
            }
            return nil
        })
        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }
        if err := mock.ExpectationsWereMet(); err != nil {
            t.Fatalf("unmet expectations: %v", err)
        }
    })

    t.Run("rollback on failure", func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO table_a").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectExec("INSERT INTO table_b").WithArgs(2).WillReturnError(errors.New("insert failed"))
        mock.ExpectRollback()

        err := database.WithTx(context.Background(), db, func(tx *sql.Tx) error {
            if _, e := tx.ExecContext(context.Background(), "INSERT INTO table_a (id) VALUES (?)", 1); e != nil {
                return e
            }
            if _, e := tx.ExecContext(context.Background(), "INSERT INTO table_b (id) VALUES (?)", 2); e != nil {
                return e
            }
            return nil
        })
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if err := mock.ExpectationsWereMet(); err != nil {
            t.Fatalf("unmet expectations: %v", err)
        }
    })
}
