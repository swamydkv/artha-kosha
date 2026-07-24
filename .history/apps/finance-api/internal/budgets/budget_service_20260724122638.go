package budgets
package budgets

import (
    "context"
    "database/sql"
    "errors"

    "artha-kosha/apps/finance-api/internal/database"
)

type CreateBudgetRequest struct {
    OwnerID string
    Amount  int64
    Name    string
}

type Budget struct {
    BudgetID string
    OwnerID  string
    Amount   int64
}

type Repository interface {
    CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error)
}

type Service struct {
    repo Repository
    db   *sql.DB
}

func NewService(repo Repository) *Service { return &Service{repo: repo, db: nil} }
func NewServiceWithDB(repo Repository, db *sql.DB) *Service { return &Service{repo: repo, db: db} }

func (s *Service) CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error) {
    if s.db != nil {
        var b *Budget
        if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
            bb, e := s.repo.CreateBudget(ctx, req)
            if e != nil {
                return e
            }
            b = bb
            return nil
        }); err != nil {
            return nil, errors.New("failed to create budget")
        }
        return b, nil
    }
    return s.repo.CreateBudget(ctx, req)
}
