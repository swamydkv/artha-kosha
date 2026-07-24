package contract

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestDeleteSessionsContract_T027(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Session Revoke All User",
		DateOfBirth:     "1990-01-15",
		MobileNumber:    "+1987654331",
		Email:           "sessionrevokeall@example.com",
		Username:        "sessionrevokeall",
		Password:        "SessionPass123!",
		ConfirmPassword: "SessionPass123!",
	}
	registerResp, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Login twice
	loginResp1, err := provider.Login(auth.LoginRequest{Username: "sessionrevokeall", Password: "SessionPass123!"})
	if err != nil {
		t.Fatalf("Login 1 failed: %v", err)
	}
	loginResp2, err := provider.Login(auth.LoginRequest{Username: "sessionrevokeall", Password: "SessionPass123!"})
	if err != nil {
		t.Fatalf("Login 2 failed: %v", err)
	}

	err = provider.RevokeAll(registerResp.UserID)
	if err != nil {
		t.Fatalf("Failed to revoke all sessions: %v", err)
	}

	// Sessions should no longer be valid
	_, err = provider.GetSession(loginResp1.SessionID)
	if err == nil {
		t.Error("Expected error when getting revoked session 1")
	}
	_, err = provider.GetSession(loginResp2.SessionID)
	if err == nil {
		t.Error("Expected error when getting revoked session 2")
	}
}
