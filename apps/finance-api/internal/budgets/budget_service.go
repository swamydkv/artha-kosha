package budgets

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

type DomainEmitter interface {
	EmitTx(ctx context.Context, tx *sql.Tx, e domain.DomainEvent) error
}

type Service struct {
	repo   Repository
	db     *sql.DB
	audit  *audit.Service
	domain DomainEmitter
}

func NewService(repo Repository) *Service                   { return &Service{repo: repo, db: nil} }
func NewServiceWithDB(repo Repository, db *sql.DB) *Service { return &Service{repo: repo, db: db} }
func (s *Service) SetAuditService(a *audit.Service)         { s.audit = a }
func (s *Service) SetDomainService(d DomainEmitter)       { s.domain = d }

func (s *Service) CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error) {
	if s.db != nil {
		var b *Budget
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			bb, e := s.repo.CreateBudget(ctx, req)
			if e != nil {
				return e
			}
			b = bb
			if s.audit != nil && b != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					UserID:     req.OwnerID,
					Resource:   "budget",
					ResourceID: b.BudgetID,
					Action:     "create",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			if s.domain != nil && b != nil {
				ev := domain.DomainEvent{
					ID:            genID("de"),
					EventType:     "BUDGET_CREATED",
					AggregateID:   b.BudgetID,
					AggregateType: "budget",
					EventData:     nil,
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return nil, errors.New("failed to create budget")
		}
		if s.audit != nil && b != nil {
			_ = s.audit.Record(ctx, audit.AuditEvent{
				UserID:     req.OwnerID,
				Resource:   "budget",
				ResourceID: b.BudgetID,
				Action:     "create",
				Result:     "success",
			})
		}
		return b, nil
	}
	return s.repo.CreateBudget(ctx, req)
}

func genID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return prefix + "-" + hex.EncodeToString(buf)
}
