package budgets

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"artha-kosha/apps/finance-api/internal/sqlc"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error) {
	q := sqlc.New(r.db)
	id := uuid.New()
	if err := q.InsertBudget(ctx, sqlc.InsertBudgetParams{ID: id, OwnerID: uuid.MustParse(req.OwnerID), Amount: req.Amount, Name: req.Name}); err != nil {
		return nil, err
	}
	return &Budget{BudgetID: id.String(), OwnerID: req.OwnerID, Amount: req.Amount, Name: req.Name}, nil
}
