package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestRegistrationConflict_DuplicateUsername(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register first user
	request1 := auth.RegisterUserRequest{
		FullName:        "First User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1111111111",
		Email:           "first@example.com",
		Username:        "sameuser",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	}

	_, err := provider.Register(request1)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register second user with same username
	request2 := auth.RegisterUserRequest{
		FullName:        "Second User",
		DateOfBirth:     "1995-01-01",
		MobileNumber:    "+2222222222",
		Email:           "second@example.com",
		Username:        "sameuser", // Duplicate username
		Password:        "Password456!",
		ConfirmPassword: "Password456!",
	}

	_, err = provider.Register(request2)
	if err == nil {
		t.Error("Expected error when registering duplicate username")
	}
}

func TestRegistrationConflict_DuplicateEmail(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register first user
	request1 := auth.RegisterUserRequest{
		FullName:        "First User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1111111111",
		Email:           "same@example.com",
		Username:        "user1",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	}

	_, err := provider.Register(request1)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register second user with same email
	request2 := auth.RegisterUserRequest{
		FullName:        "Second User",
		DateOfBirth:     "1995-01-01",
		MobileNumber:    "+2222222222",
		Email:           "same@example.com", // Duplicate email
		Username:        "user2",
		Password:        "Password456!",
		ConfirmPassword: "Password456!",
	}

	_, err = provider.Register(request2)
	if err == nil {
		t.Error("Expected error when registering duplicate email")
	}
}

func TestRegistrationConflict_DuplicateMobileNumber(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register first user
	request1 := auth.RegisterUserRequest{
		FullName:        "First User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "first@example.com",
		Username:        "user1",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	}

	_, err := provider.Register(request1)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register second user with same mobile number
	request2 := auth.RegisterUserRequest{
		FullName:        "Second User",
		DateOfBirth:     "1995-01-01",
		MobileNumber:    "+1234567890", // Duplicate mobile number
		Email:           "second@example.com",
		Username:        "user2",
		Password:        "Password456!",
		ConfirmPassword: "Password456!",
	}

	_, err = provider.Register(request2)
	if err == nil {
		t.Error("Expected error when registering duplicate mobile number")
	}
}
