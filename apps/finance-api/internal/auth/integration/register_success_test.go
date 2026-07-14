package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestRegistrationIntegration_Success(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	request := auth.RegisterUserRequest{
		FullName:        "Integration User",
		DateOfBirth:     "1992-03-10",
		MobileNumber:    "+1555123456",
		Email:           "integration@example.com",
		Username:        "integrationuser",
		Password:        "IntegrationPass123!",
		ConfirmPassword: "IntegrationPass123!",
	}

	response, err := provider.Register(request)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Verify user was created successfully
	if response.UserID == "" {
		t.Error("Expected non-empty UserID")
	}
	if response.Username != "integrationuser" {
		t.Errorf("Expected username 'integrationuser', got '%s'", response.Username)
	}
	if response.FirstName != "Integration" {
		t.Errorf("Expected first name 'Integration', got '%s'", response.FirstName)
	}
	if response.PasswordHash == "" {
		t.Error("Expected password hash to be generated")
	}

	// Verify duplicate registration fails
	_, err = provider.Register(request)
	if err == nil {
		t.Error("Expected error when registering duplicate user")
	}
}

func TestRegistrationIntegration_Persistence(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	request := auth.RegisterUserRequest{
		FullName:        "Persistence User",
		DateOfBirth:     "1988-11-25",
		MobileNumber:    "+1444555666",
		Email:           "persistence@example.com",
		Username:        "persistenceuser",
		Password:        "PersistencePass456!",
		ConfirmPassword: "PersistencePass456!",
	}

	response, err := provider.Register(request)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to login with the registered user
	loginRequest := auth.LoginRequest{
		Username: "persistenceuser",
		Password: "PersistencePass456!",
	}

	loginResponse, err := provider.Login(loginRequest)
	if err != nil {
		t.Errorf("Login failed after registration: %v", err)
	}

	if loginResponse.UserID != response.UserID {
		t.Errorf("UserID mismatch: expected %s, got %s", response.UserID, loginResponse.UserID)
	}

	if loginResponse.WelcomeMessage == "" {
		t.Error("Expected welcome message")
	}
}