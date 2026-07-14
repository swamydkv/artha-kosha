package contract

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestLoginContract(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// First register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Login Test User",
		DateOfBirth:     "1990-01-15",
		MobileNumber:    "+1987654321",
		Email:           "logintest@example.com",
		Username:        "logintestuser",
		Password:        "LoginPass123!",
		ConfirmPassword: "LoginPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name        string
		request     auth.LoginRequest
		expectError bool
	}{
		{
			name: "Valid login",
			request: auth.LoginRequest{
				Username: "logintestuser",
				Password: "LoginPass123!",
			},
			expectError: false,
		},
		{
			name: "Invalid username",
			request: auth.LoginRequest{
				Username: "nonexistentuser",
				Password: "LoginPass123!",
			},
			expectError: true,
		},
		{
			name: "Invalid password",
			request: auth.LoginRequest{
				Username: "logintestuser",
				Password: "WrongPassword!",
			},
			expectError: true,
		},
		{
			name: "Empty username",
			request: auth.LoginRequest{
				Username: "",
				Password: "LoginPass123!",
			},
			expectError: true,
		},
		{
			name: "Empty password",
			request: auth.LoginRequest{
				Username: "logintestuser",
				Password: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Login(tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestLoginResponseContract(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Response Test User",
		DateOfBirth:     "1985-05-20",
		MobileNumber:    "+1555555555",
		Email:           "responsetest@example.com",
		Username:        "responsetestuser",
		Password:        "ResponsePass456!",
		ConfirmPassword: "ResponsePass456!",
	}

	registerResponse, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Login
	loginRequest := auth.LoginRequest{
		Username: "responsetestuser",
		Password: "ResponsePass456!",
	}

	loginResponse, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Verify response structure
	if loginResponse.UserID == "" {
		t.Error("UserID should not be empty")
	}
	if loginResponse.SessionID == "" {
		t.Error("SessionID should not be empty")
	}
	if loginResponse.WelcomeMessage == "" {
		t.Error("WelcomeMessage should not be empty")
	}
	if loginResponse.UserID != registerResponse.UserID {
		t.Errorf("UserID mismatch: expected %s, got %s", registerResponse.UserID, loginResponse.UserID)
	}
}