package validation

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestLoginValidation_RequiredFields(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a test user first
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Validation Test User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1999888777",
		Email:           "validation@example.com",
		Username:        "validationuser",
		Password:        "ValidationPass123!",
		ConfirmPassword: "ValidationPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
	}{
		{
			name:        "Empty username",
			username:    "",
			password:    "ValidationPass123!",
			expectError: true,
		},
		{
			name:        "Empty password",
			username:    "validationuser",
			password:    "",
			expectError: true,
		},
		{
			name:        "Both empty",
			username:    "",
			password:    "",
			expectError: true,
		},
		{
			name:        "Valid credentials",
			username:    "validationuser",
			password:    "ValidationPass123!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginRequest := auth.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			_, err := provider.Login(loginRequest)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestLoginValidation_CredentialMatching(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a test user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Credential Test User",
		DateOfBirth:     "1985-05-20",
		MobileNumber:    "+122211112222",
		Email:           "credential@example.com",
		Username:        "credentialuser",
		Password:        "CredentialPass123!",
		ConfirmPassword: "CredentialPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
	}{
		{
			name:        "Correct username, wrong password",
			username:    "credentialuser",
			password:    "WrongPassword!",
			expectError: true,
		},
		{
			name:        "Wrong username, correct password",
			username:    "wronguser",
			password:    "CredentialPass123!",
			expectError: true,
		},
		{
			name:        "Both wrong",
			username:    "wronguser",
			password:    "WrongPassword!",
			expectError: true,
		},
		{
			name:        "Case insensitive username",
			username:    "CredentialUser", // Different case but should still work
			password:    "CredentialPass123!",
			expectError: false, // Username lookup is case-insensitive
		},
		{
			name:        "Case sensitive password",
			username:    "credentialuser",
			password:    "credentialpass123!", // Different case
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginRequest := auth.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			_, err := provider.Login(loginRequest)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
