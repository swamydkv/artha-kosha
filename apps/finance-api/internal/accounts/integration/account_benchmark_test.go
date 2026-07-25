package integration

import (
	"context"
	"testing"

	"artha-kosha/apps/finance-api/internal/accounts"
)

type mockAccountRepo struct{}

func (m mockAccountRepo) CreateAccount(ctx context.Context, req accounts.CreateAccountRequest) (*accounts.Account, error) {
	return &accounts.Account{AccountID: "acc-1", OwnerID: req.OwnerID, Name: req.Name}, nil
}

func BenchmarkCreateAccount(b *testing.B) {
	svc := accounts.NewService(mockAccountRepo{})
	ctx := context.Background()
	req := accounts.CreateAccountRequest{OwnerID: "user-1", Name: "Savings"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.CreateAccount(ctx, req)
	}
}
