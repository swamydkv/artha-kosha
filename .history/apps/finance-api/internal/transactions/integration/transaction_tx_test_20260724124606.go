package integration

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/transactions"
)

func TestTransactionService_CreateTransaction_Transaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO audit_events").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &fakeTxRepo{}
	svc := transactions.NewServiceWithDB(repo, db)
	svc.SetAuditService(audit.NewService(audit.NewSQLRepository(db)))
	_, err = svc.CreateTransaction(context.Background(), transactions.CreateTransactionRequest{AccountID: "acc-1", Amount: 100})
	if err != nil {
		t.Fatalf("create transaction failed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

type fakeTxRepo struct{}

func (f *fakeTxRepo) CreateTransaction(ctx context.Context, req transactions.CreateTransactionRequest) (*transactions.Transaction, error) {
	return &transactions.Transaction{TransactionID: "tx-1", AccountID: req.AccountID, Amount: req.Amount}, nil
}
