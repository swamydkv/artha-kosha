package integration

import (
	"context"
	"testing"

	"artha-kosha/apps/finance-api/internal/transactions"
)

type mockTransactionRepo struct{}

func (m mockTransactionRepo) CreateTransaction(ctx context.Context, req transactions.CreateTransactionRequest) (*transactions.Transaction, error) {
	return &transactions.Transaction{TransactionID: "tx-1", AccountID: req.AccountID, Amount: req.Amount}, nil
}

func BenchmarkCreateTransaction(b *testing.B) {
	svc := transactions.NewService(mockTransactionRepo{})
	ctx := context.Background()
	req := transactions.CreateTransactionRequest{AccountID: "acc-1", Amount: 100, Memo: "Grocery"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.CreateTransaction(ctx, req)
	}
}
