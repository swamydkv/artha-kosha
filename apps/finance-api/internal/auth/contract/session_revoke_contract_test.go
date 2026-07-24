package contract

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestDeleteSessionContract_T026(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Session Revoke User",
		DateOfBirth:     "1990-01-15",
		MobileNumber:    "+1987654330",
		Email:           "sessionrevokec@example.com",
		Username:        "sessionrevokec",
		Password:        "SessionPass123!",
		ConfirmPassword: "SessionPass123!",
	}
	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Login
	loginResp, err := provider.Login(auth.LoginRequest{Username: "sessionrevokec", Password: "SessionPass123!"})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Revoke
	err = provider.Logout(loginResp.SessionID)
	if err != nil {
		t.Fatalf("Failed to revoke session: %v", err)
	}

	// Session should no longer be valid
	_, err = provider.GetSession(loginResp.SessionID)
	if err == nil {
		t.Error("Expected error when getting revoked session")
	}
}
