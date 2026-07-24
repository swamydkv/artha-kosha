package transactions
package transactions

import (
	"context"
	"database/sql"
	"errors"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/database"
)

type CreateTransactionRequest struct {
	AccountID string
	Amount    int64
	Memo      string
}

type Transaction struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

type Repository interface {
	CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error)
}

type Service struct {
	repo  Repository
	db    *sql.DB
	audit *audit.Service
}

func NewService(repo Repository) *Service { return &Service{repo: repo, db: nil} }
func NewServiceWithDB(repo Repository, db *sql.DB) *Service { return &Service{repo: repo, db: db} }

func (s *Service) SetAuditService(a *audit.Service) { s.audit = a }

func (s *Service) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	if s.db != nil {
		var tr *Transaction
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			t, e := s.repo.CreateTransaction(ctx, req)
			if e != nil {
				return e
			}
			tr = t
			return nil
		}); err != nil {
			return nil, errors.New("failed to create transaction")
		}
		if s.audit != nil && tr != nil {
			_ = s.audit.Record(ctx, audit.AuditEvent{
				UserID:     "",
				Resource:   "transaction",
				ResourceID: tr.TransactionID,
				Action:     "create",
				Result:     "success",
			})
		}
		return tr, nil
	}
	return s.repo.CreateTransaction(ctx, req)
}
