package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestAuthFlow_CompleteLifecycle(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Flow Test User",
		DateOfBirth:     "1992-03-10",
		MobileNumber:    "+1333444555",
		Email:           "flowtest@example.com",
		Username:        "flowtestuser",
		Password:        "FlowPass123!",
		ConfirmPassword: "FlowPass123!",
	}

	registerResponse, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Login with the registered user
	loginRequest := auth.LoginRequest{
		Username: "flowtestuser",
		Password: "FlowPass123!",
	}

	loginResponse, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Verify login response
	if loginResponse.UserID != registerResponse.UserID {
		t.Errorf("UserID mismatch: expected %s, got %s", registerResponse.UserID, loginResponse.UserID)
	}
	if loginResponse.SessionID == "" {
		t.Error("SessionID should not be empty")
	}
	if loginResponse.WelcomeMessage == "" {
		t.Error("WelcomeMessage should not be empty")
	}

	// Logout
	err = provider.Logout(loginResponse.SessionID)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Verify session is invalid after logout
	// This would require session validation, which is not implemented in the current provider
	// For now, we just verify that logout doesn't error
}

func TestAuthFlow_MultipleLogins(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Multi Login User",
		DateOfBirth:     "1988-11-25",
		MobileNumber:    "+1444666777",
		Email:           "multilogin@example.com",
		Username:        "multiloginuser",
		Password:        "MultiLoginPass456!",
		ConfirmPassword: "MultiLoginPass456!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// First login
	loginRequest := auth.LoginRequest{
		Username: "multiloginuser",
		Password: "MultiLoginPass456!",
	}

	session1, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("First login failed: %v", err)
	}

	// Second login (should create a new session)
	session2, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("Second login failed: %v", err)
	}

	// Verify different sessions were created
	if session1.SessionID == session2.SessionID {
		t.Error("Each login should create a new session")
	}

	// Verify same user
	if session1.UserID != session2.UserID {
		t.Error("UserIDs should match for same user")
	}
}