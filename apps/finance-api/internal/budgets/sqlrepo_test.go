package budgets

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSQLRepository_CreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	ownerID := uuid.New().String()
	mock.ExpectExec("INSERT INTO budgets").WithArgs(sqlmock.AnyArg(), ownerID, sqlmock.AnyArg(), "monthly").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLRepository(db)
	req := CreateBudgetRequest{OwnerID: ownerID, Amount: 500, Name: "monthly"}
	b, err := repo.CreateBudget(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil || b.OwnerID != ownerID {
		t.Fatalf("unexpected budget result: %#v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
