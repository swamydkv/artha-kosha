package integration

import (
	"context"
	"testing"

	"artha-kosha/apps/finance-api/internal/budgets"
)

type mockBudgetRepo struct{}

func (m mockBudgetRepo) CreateBudget(ctx context.Context, req budgets.CreateBudgetRequest) (*budgets.Budget, error) {
	return &budgets.Budget{BudgetID: "bdg-1", OwnerID: req.OwnerID, Amount: req.Amount}, nil
}

func BenchmarkCreateBudget(b *testing.B) {
	svc := budgets.NewService(mockBudgetRepo{})
	ctx := context.Background()
	req := budgets.CreateBudgetRequest{OwnerID: "user-1", Amount: 1000, Name: "Food"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.CreateBudget(ctx, req)
	}
}
