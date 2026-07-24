package transactions

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/google/uuid"

	"artha-kosha/apps/finance-api/internal/sqlc"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	q := sqlc.New(r.db)
	id := uuid.New()
	memo := sql.NullString{String: req.Memo, Valid: req.Memo != ""}
	if err := q.InsertTransaction(ctx, sqlc.InsertTransactionParams{ID: id, AccountID: uuid.MustParse(req.AccountID), Amount: strconv.FormatInt(req.Amount, 10), Memo: memo}); err != nil {
		return nil, err
	}
	return &Transaction{TransactionID: id.String(), AccountID: req.AccountID, Amount: req.Amount}, nil
}
