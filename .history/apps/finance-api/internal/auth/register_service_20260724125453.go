package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/database"
	"artha-kosha/apps/finance-api/internal/domain"
	"artha-kosha/apps/finance-api/internal/users"
)

// RegisterService handles user registration logic
type RegisterService struct {
	authProvider AuthProvider
	userRepo     users.Repository
	db           *sql.DB // optional DB for transactional operations
	audit        *audit.Service
	domain       *domain.Service
}

// NewRegisterService creates a new registration service (non-DB-backed)
func NewRegisterService(authProvider AuthProvider, userRepo users.Repository) *RegisterService {
	return &RegisterService{
		authProvider: authProvider,
		userRepo:     userRepo,
		db:           nil,
	}
}

// NewRegisterServiceWithDB creates a DB-backed registration service that uses transactions
func NewRegisterServiceWithDB(authProvider AuthProvider, userRepo users.Repository, db *sql.DB) *RegisterService {
	return &RegisterService{
		authProvider: authProvider,
		userRepo:     userRepo,
		db:           db,
	}
}

func (s *RegisterService) SetAuditService(a *audit.Service) { s.audit = a }

func (s *RegisterService) SetDomainService(d *domain.Service) { s.domain = d }

// Register handles the complete registration flow
func (s *RegisterService) Register(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error) {
	// Validate the request
	if err := validateRegistrationRequest(req); err != nil {
		return nil, err
	}

	// Check for existing users
	exists, err := s.userRepo.CheckUserExists(ctx, req.Username, req.Email, req.MobileNumber)
	if err != nil {
		return nil, errors.New("failed to check existing users")
	}

	if exists.UsernameExists {
		return nil, errors.New("username already exists")
	}
	if exists.EmailExists {
		return nil, errors.New("email already exists")
	}
	if exists.MobileExists {
		return nil, errors.New("mobile number already exists")
	}

	// Parse date of birth
	dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date of birth format")
	}

	// Hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user request
	createReq := users.CreateUserRequest{
		FullName:     req.FullName,
		DateOfBirth:  dateOfBirth,
		MobileNumber: req.MobileNumber,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	var created *users.User

	// If DB available, run creation in a transaction
	if s.db != nil {
		if err := database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
			u, e := s.userRepo.CreateUser(ctx, createReq)
			if e != nil {
				return e
			}
			created = u
			if s.audit != nil && created != nil {
				if err := s.audit.RecordTx(ctx, tx, audit.AuditEvent{
					UserID:     created.UserID,
					Resource:   "user",
					ResourceID: created.UserID,
					Action:     "register",
					Result:     "success",
				}); err != nil {
					return err
				}
			}
			// emit domain event for user registration inside the same tx
			if s.domain != nil && created != nil {
				ev := domain.DomainEvent{
					ID:            genID("de"),
					EventType:     "USER_REGISTERED",
					AggregateID:   created.UserID,
					AggregateType: "user",
					EventData:     nil,
				}
				if err := s.domain.EmitTx(ctx, tx, ev); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return nil, errors.New("failed to create user")
		}
	} else {
		u, err := s.userRepo.CreateUser(ctx, createReq)
		if err != nil {
			return nil, errors.New("failed to create user")
		}
		created = u
	}

	// Generate response
	firstName := firstNameFromFullName(req.FullName)
	// emit audit event (best-effort)
	if s.audit != nil && created != nil {
		_ = s.audit.Record(ctx, audit.AuditEvent{
			UserID:     created.UserID,
			Resource:   "user",
			ResourceID: created.UserID,
			Action:     "register",
			Result:     "success",
		})
	}
	return &RegisterUserResponse{
		UserID:       created.UserID,
		Username:     created.Username,
		PasswordHash: created.PasswordHash,
		FirstName:    firstName,
	}, nil
}
