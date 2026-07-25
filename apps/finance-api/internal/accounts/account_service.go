package accounts

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

type CreateAccountRequest struct {
	OwnerID string
	Name    string
}

type Account struct {
	AccountID string
	OwnerID   string
	Name      string
}

type Repository interface {
	CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error)
}

type DomainEmitter interface {
	EmitTx(ctx context.Context, tx *sql.Tx, e domain.DomainEvent) error
}

type Service struct {
	repo   Repository
	db     *sql.DB
	audit  *audit.Service
	domain DomainEmitter
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, db: nil}
}

func NewServiceWithDB(repo Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) SetAuditService(a *audit.Service)   { s.audit = a }
func (s *Service) SetDomainService(d DomainEmitter) { s.domain = d }

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	if s.db != nil {
		var acc *Account
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			a, e := s.repo.CreateAccount(ctx, req)
			if e != nil {
				return e
			}
			acc = a
			// insert audit within the same tx if available
			if s.audit != nil && acc != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					UserID:     req.OwnerID,
					Resource:   "account",
					ResourceID: acc.AccountID,
					Action:     "create",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			if s.domain != nil && acc != nil {
				ev := domain.DomainEvent{
					ID:            genID("de"),
					EventType:     "ACCOUNT_CREATED",
					AggregateID:   acc.AccountID,
					AggregateType: "account",
					EventData:     nil,
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return nil, errors.New("failed to create account")
		}
		if s.audit != nil && acc != nil {
			_ = s.audit.Record(ctx, audit.AuditEvent{
				UserID:     req.OwnerID,
				Resource:   "account",
				ResourceID: acc.AccountID,
				Action:     "create",
				Result:     "success",
			})
		}
		return acc, nil
	}
	return s.repo.CreateAccount(ctx, req)
}

func genID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return prefix + "-" + hex.EncodeToString(buf)
}
