package auth

import (
	"context"
	"errors"

	"artha-kosha/apps/finance-api/internal/users"
)

// LoginService handles user authentication logic
type LoginService struct {
	authProvider AuthProvider
	userRepo     users.Repository
}

// NewLoginService creates a new login service
func NewLoginService(authProvider AuthProvider, userRepo users.Repository) *LoginService {
	return &LoginService{
		authProvider: authProvider,
		userRepo:     userRepo,
	}
}

// Login handles the complete login flow
func (s *LoginService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validate the request
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
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

	return s.authProvider.Logout(sessionID)
}