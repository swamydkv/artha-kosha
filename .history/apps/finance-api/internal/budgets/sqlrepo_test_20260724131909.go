package budgets

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLRepository_CreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO budgets").WithArgs(sqlmock.AnyArg(), "owner1", sqlmock.AnyArg(), "monthly").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLRepository(db)
	req := CreateBudgetRequest{OwnerID: "owner1", Amount: 500.0, Name: "monthly"}
	b, err := repo.CreateBudget(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil || b.OwnerID != "owner1" {
		t.Fatalf("unexpected budget result: %#v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
