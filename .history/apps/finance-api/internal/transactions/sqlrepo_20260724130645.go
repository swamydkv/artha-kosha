package transactions

import (
	"context"
	"database/sql"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	q := `INSERT INTO transactions (id, account_id, amount, memo) VALUES ($1,$2,$3,$4)`
	id := req.AccountID + "-tx"
	if _, err := r.db.ExecContext(ctx, q, id, req.AccountID, req.Amount, req.Memo); err != nil {
		return nil, err
	}
	return &Transaction{TransactionID: id, AccountID: req.AccountID, Amount: req.Amount}, nil
}
