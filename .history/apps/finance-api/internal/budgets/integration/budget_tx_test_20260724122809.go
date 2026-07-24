package integration
package integration

import (
    "context"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"

    "artha-kosha/apps/finance-api/internal/budgets"
)

func TestBudgetService_CreateBudget_Transaction(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock new: %v", err)
    }
    defer db.Close()

    mock.ExpectBegin()
    mock.ExpectCommit()

    repo := &fakeBudgetRepo{}
    svc := budgets.NewServiceWithDB(repo, db)
    _, err = svc.CreateBudget(context.Background(), budgets.CreateBudgetRequest{OwnerID: "owner1", Amount: 500, Name: "Monthly"})
    if err != nil {
        t.Fatalf("create budget failed: %v", err)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

type fakeBudgetRepo struct{}
func (f *fakeBudgetRepo) CreateBudget(ctx context.Context, req budgets.CreateBudgetRequest) (*budgets.Budget, error) {
    return &budgets.Budget{BudgetID: "b-1", OwnerID: req.OwnerID, Amount: req.Amount}, nil
}
