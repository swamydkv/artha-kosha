package auth

import (
	"context"
	"database/sql"
	"errors"
	import_log "log"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/database"
	"artha-kosha/apps/finance-api/internal/domain"
	"github.com/google/uuid"
	"artha-kosha/apps/finance-api/internal/sessions"
	"artha-kosha/apps/finance-api/internal/transactions"
	"artha-kosha/apps/finance-api/internal/users"
)

// LoginService handles user authentication logic
type LoginService struct {
	authProvider AuthProvider
	userRepo     users.Repository
	db           *sql.DB // optional DB for transactional operations
	audit        *audit.Service
	domain       *domain.Service
}

// NewLoginService creates a new login service (non-DB-backed)
func NewLoginService(authProvider AuthProvider, userRepo users.Repository) *LoginService {
	return &LoginService{
		authProvider: authProvider,
		userRepo:     userRepo,
		db:           nil,
	}
}

// NewLoginServiceWithDB creates a DB-backed login service that uses transactions
func NewLoginServiceWithDB(authProvider AuthProvider, userRepo users.Repository, db *sql.DB) *LoginService {
	return &LoginService{
		authProvider: authProvider,
		userRepo:     userRepo,
		db:           db,
	}
}

func (s *LoginService) SetAuditService(a *audit.Service) { s.audit = a }

func (s *LoginService) SetDomainService(d *domain.Service) { s.domain = d }

// LoginCtx handles the complete login flow with context
func (s *LoginService) LoginCtx(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validate the request
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// If DB available, run provider login inside a transaction
	if s.db != nil {
		var resp LoginResponse
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			r, e := s.authProvider.Login(req)
			if e != nil {
				return e
			}
			resp = r
			if s.audit != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					ID:         uuid.New().String(),
					UserID:     resp.UserID,
					Resource:   "session",
					ResourceID: resp.SessionID,
					Action:     "login",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			if s.domain != nil {
				ev := domain.DomainEvent{
					ID:            uuid.New().String(),
					EventType:     "USER_LOGGED_IN",
					AggregateID:   resp.UserID,
					AggregateType: "user",
					EventData:     []byte("{}"),
					ProcessingStatus: "pending",
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			import_log.Printf("Login transaction failed: %v", err)
			return nil, errors.New("invalid credentials")
		}
		return &resp, nil
	}

	// Use the auth provider to authenticate
	response, err := s.authProvider.Login(req)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &response, nil
}



func (s *LoginService) Login(req LoginRequest) (LoginResponse, error) {
	resp, err := s.LoginCtx(context.Background(), req)
	if err != nil {
		return LoginResponse{}, err
	}
	return *resp, nil
}

// LogoutCtx handles the logout flow with context
func (s *LoginService) LogoutCtx(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	if s.db != nil {
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			if err := s.authProvider.Logout(sessionID); err != nil {
				return err
			}
			if s.audit != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					ID:         uuid.New().String(),
					Resource:   "session",
					ResourceID: sessionID,
					Action:     "logout",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			if s.domain != nil {
				ev := domain.DomainEvent{
					ID:            uuid.New().String(),
					EventType:     "USER_LOGGED_OUT",
					AggregateID:   sessionID,
					AggregateType: "session",
					EventData:     []byte("{}"),
					ProcessingStatus: "pending",
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	if err := s.authProvider.Logout(sessionID); err != nil {
		return err
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, audit.AuditEvent{
			ID:         uuid.New().String(),
			Resource:   "session",
			ResourceID: sessionID,
			Action:     "logout",
			Result:     "success",
		})
	}
	return nil
}

func (s *LoginService) Logout(sessionID string) error {
	return s.LogoutCtx(context.Background(), sessionID)
}

func (s *LoginService) Register(req RegisterUserRequest) (RegisterUserResponse, error) {
	return s.authProvider.Register(req)
}

func (s *LoginService) GetSession(id string) (sessions.Session, error) {
	return s.authProvider.GetSession(id)
}

func (s *LoginService) RevokeAll(id string) error {
	return s.authProvider.RevokeAll(id)
}

func (s *LoginService) ChangePassword(req ChangePasswordRequest) error {
	return s.authProvider.ChangePassword(req)
}

func (s *LoginService) DeleteUser(userID string, confirmation string, archiveRetentionDays int) error {
	return s.authProvider.DeleteUser(userID, confirmation, archiveRetentionDays)
}

func (s *LoginService) GetAccountsService() *accounts.Service {
	if p, ok := s.authProvider.(interface{ GetAccountsService() *accounts.Service }); ok {
		return p.GetAccountsService()
	}
	return nil
}

func (s *LoginService) GetTransactionsService() *transactions.Service {
	if p, ok := s.authProvider.(interface{ GetTransactionsService() *transactions.Service }); ok {
		return p.GetTransactionsService()
	}
	return nil
}

func (s *LoginService) GetBudgetsService() *budgets.Service {
	if p, ok := s.authProvider.(interface{ GetBudgetsService() *budgets.Service }); ok {
		return p.GetBudgetsService()
	}
	return nil
}
