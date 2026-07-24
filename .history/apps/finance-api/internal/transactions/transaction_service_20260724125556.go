package transactions

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/database"
	"artha-kosha/apps/finance-api/internal/domain"
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
	domain *domain.Service
}

func NewService(repo Repository) *Service                   { return &Service{repo: repo, db: nil} }
func NewServiceWithDB(repo Repository, db *sql.DB) *Service { return &Service{repo: repo, db: db} }

func (s *Service) SetAuditService(a *audit.Service) { s.audit = a }
func (s *Service) SetDomainService(d *domain.Service) { s.domain = d }

func (s *Service) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	if s.db != nil {
		var tr *Transaction
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			t, e := s.repo.CreateTransaction(ctx, req)
			if e != nil {
				return e
			}
			tr = t
			if s.audit != nil && tr != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					UserID:     "",
					Resource:   "transaction",
					ResourceID: tr.TransactionID,
					Action:     "create",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			if s.domain != nil && tr != nil {
				ev := domain.DomainEvent{
					ID:            genID("de"),
					EventType:     "TRANSACTION_CREATED",
					AggregateID:   tr.TransactionID,
					AggregateType: "transaction",
					EventData:     nil,
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
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
