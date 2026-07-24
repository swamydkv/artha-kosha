package transactions

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLRepository_CreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO transactions").WithArgs(sqlmock.AnyArg(), "acc1", sqlmock.AnyArg(), "note").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLRepository(db)
	req := CreateTransactionRequest{AccountID: "acc1", Amount: 100.0, Memo: "note"}
	tx, err := repo.CreateTransaction(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil || tx.AccountID != "acc1" {
		t.Fatalf("unexpected tx result: %#v", tx)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
