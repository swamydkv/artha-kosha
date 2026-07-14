package users

import (
	"time"
)

// User represents a user account in the system
type User struct {
	UserID       string
	FullName     string
	DateOfBirth  time.Time
	MobileNumber string
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	FullName        string
	DateOfBirth     time.Time
	MobileNumber    string
	Email           string
	Username        string
	PasswordHash    string
}

// UserExistsCheck represents the result of checking if a user exists
type UserExistsCheck struct {
	UsernameExists bool
	EmailExists     bool
	MobileExists    bool
}