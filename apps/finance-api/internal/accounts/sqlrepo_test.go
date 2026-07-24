package accounts

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSQLRepository_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	ownerID := uuid.New().String()
	mock.ExpectExec("INSERT INTO accounts").WithArgs(sqlmock.AnyArg(), ownerID, "My Account").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLRepository(db)
	req := CreateAccountRequest{OwnerID: ownerID, Name: "My Account"}
	acc, err := repo.CreateAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc == nil || acc.OwnerID != ownerID {
		t.Fatalf("unexpected account result: %#v", acc)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
