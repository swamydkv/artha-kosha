package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"artha-kosha/apps/finance-api/internal/users"
)

// RegisterService handles user registration logic
type RegisterService struct {
	authProvider AuthProvider
	userRepo     users.Repository
}

// NewRegisterService creates a new registration service
func NewRegisterService(authProvider AuthProvider, userRepo users.Repository) *RegisterService {
	return &RegisterService{
		authProvider: authProvider,
		userRepo:     userRepo,
	}
}

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

	// Create user in repository
	createReq := users.CreateUserRequest{
		FullName:     req.FullName,
		DateOfBirth:  dateOfBirth,
		MobileNumber: req.MobileNumber,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	user, err := s.userRepo.CreateUser(ctx, createReq)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate response
	firstName := firstNameFromFullName(req.FullName)
	return &RegisterUserResponse{
		UserID:       user.UserID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		FirstName:    firstName,
	}, nil
}
