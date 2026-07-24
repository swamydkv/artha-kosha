package users

import (
	"context"
)

// Repository defines the interface for user data access
type Repository interface {
	// CreateUser creates a new user account
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// GetUserByMobileNumber retrieves a user by their mobile number
	GetUserByMobileNumber(ctx context.Context, mobileNumber string) (*User, error)

	// CheckUserExists checks if a user exists with the given username, email, or mobile number
	CheckUserExists(ctx context.Context, username, email, mobileNumber string) (*UserExistsCheck, error)
}
