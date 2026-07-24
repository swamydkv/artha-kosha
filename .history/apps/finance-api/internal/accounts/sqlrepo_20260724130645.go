package accounts

import (
	"context"
	"database/sql"
)

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	q := `INSERT INTO accounts (id, owner_id, name) VALUES ($1,$2,$3)`
	if _, err := r.db.ExecContext(ctx, q, req.OwnerID+"-acct", req.OwnerID, req.Name); err != nil {
		return nil, err
	}
	return &Account{AccountID: req.OwnerID + "-acct", OwnerID: req.OwnerID, Name: req.Name}, nil
}
