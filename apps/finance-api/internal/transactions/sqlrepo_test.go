package transactions

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSQLRepository_CreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	accID := uuid.New().String()
	mock.ExpectExec("INSERT INTO transactions").WithArgs(sqlmock.AnyArg(), accID, sqlmock.AnyArg(), "note").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLRepository(db)
	req := CreateTransactionRequest{AccountID: accID, Amount: 100, Memo: "note"}
	tx, err := repo.CreateTransaction(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil || tx.AccountID != accID {
		t.Fatalf("unexpected tx result: %#v", tx)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
