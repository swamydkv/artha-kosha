package accounts
package accounts

import (
    "context"
    "database/sql"
    "errors"

    "artha-kosha/apps/finance-api/internal/database"
)
        "artha-kosha/apps/finance-api/internal/audit"

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

type Service struct {
    repo Repository
    db   *sql.DB
}
        audit *audit.Service

func NewService(repo Repository) *Service {
    return &Service{repo: repo, db: nil}
}

func NewServiceWithDB(repo Repository, db *sql.DB) *Service {
    return &Service{repo: repo, db: db}
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
    if s.db != nil {
        var acc *Account
        if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
            a, e := s.repo.CreateAccount(ctx, req)
            if e != nil {
                return e
            }
            acc = a
            return nil
        }); err != nil {
            return nil, errors.New("failed to create account")
        }
        return acc, nil
    }
    return s.repo.CreateAccount(ctx, req)
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
