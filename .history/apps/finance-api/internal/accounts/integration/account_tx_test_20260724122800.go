package integration

import (
    "context"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"

    "artha-kosha/apps/finance-api/internal/accounts"
)

func TestAccountService_CreateAccount_Transaction(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock new: %v", err)
    }
    defer db.Close()

    mock.ExpectBegin()
    mock.ExpectCommit()

    // repo implementation that calls Exec is not present; use a simple repo that delegates
    repo := &fakeAccountRepo{}
    svc := accounts.NewServiceWithDB(repo, db)

    _, err = svc.CreateAccount(context.Background(), accounts.CreateAccountRequest{OwnerID: "owner1", Name: "My Account"})
    if err != nil {
        t.Fatalf("create account failed: %v", err)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

type fakeAccountRepo struct{}

func (f *fakeAccountRepo) CreateAccount(ctx context.Context, req accounts.CreateAccountRequest) (*accounts.Account, error) {
    return &accounts.Account{AccountID: "acc-1", OwnerID: req.OwnerID, Name: req.Name}, nil
}
