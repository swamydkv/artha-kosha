package integration
package integration

import (
    "context"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"

    "artha-kosha/apps/finance-api/internal/audit"
)

func TestSQLRepository_Insert(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock new: %v", err)
    }
    defer db.Close()

    mock.ExpectExec("INSERT INTO audit_events").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

    repo := audit.NewSQLRepository(db)
    e := audit.AuditEvent{ID: "id-1", RequestID: "r-1", Action: "test", Result: "ok"}
    if err := repo.Insert(context.Background(), e); err != nil {
        t.Fatalf("insert failed: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}
