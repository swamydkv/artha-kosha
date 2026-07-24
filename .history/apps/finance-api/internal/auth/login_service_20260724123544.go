package auth

import (
	"context"
	"database/sql"
	"errors"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/database"
	"artha-kosha/apps/finance-api/internal/users"
)

// LoginService handles user authentication logic
type LoginService struct {
	authProvider AuthProvider
	userRepo     users.Repository
	db           *sql.DB // optional DB for transactional operations
	audit        *audit.Service
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

// Login handles the complete login flow
func (s *LoginService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
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
			return nil
		}); err != nil {
			return nil, errors.New("invalid credentials")
		}
		// audit
		if s.audit != nil {
			_ = s.audit.Record(ctx, audit.AuditEvent{
				UserID:     resp.UserID,
				Resource:   "session",
				ResourceID: resp.SessionID,
				Action:     "login",
				Result:     "success",
			})
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

// Logout handles the logout flow
func (s *LoginService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	if s.db != nil {
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			return s.authProvider.Logout(sessionID)
		}); err != nil {
			return err
		}
		if s.audit != nil {
			_ = s.audit.Record(ctx, audit.AuditEvent{
				Resource:   "session",
				ResourceID: sessionID,
				Action:     "logout",
				Result:     "success",
			})
		}
		return nil
	}

	if err := s.authProvider.Logout(sessionID); err != nil {
		return err
	}
	if s.audit != nil {
		_ = s.audit.Record(ctx, audit.AuditEvent{
			Resource:   "session",
			ResourceID: sessionID,
			Action:     "logout",
			Result:     "success",
		})
	}
	return nil
}
