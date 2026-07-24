package budgets

import (
	"context"
	"database/sql"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error) {
	q := `INSERT INTO budgets (id, owner_id, amount, name) VALUES ($1,$2,$3,$4)`
	id := req.OwnerID + "-bgt"
	if _, err := r.db.ExecContext(ctx, q, id, req.OwnerID, req.Amount, req.Name); err != nil {
		return nil, err
	}
	return &Budget{BudgetID: id, OwnerID: req.OwnerID, Amount: req.Amount}, nil
}
