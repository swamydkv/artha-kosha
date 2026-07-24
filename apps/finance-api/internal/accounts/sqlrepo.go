package accounts

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"artha-kosha/apps/finance-api/internal/sqlc"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	q := sqlc.New(r.db)
	id := uuid.New()
	err := q.InsertAccount(ctx, sqlc.InsertAccountParams{ID: id, OwnerID: uuid.MustParse(req.OwnerID), Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &Account{AccountID: id.String(), OwnerID: req.OwnerID, Name: req.Name}, nil
}
