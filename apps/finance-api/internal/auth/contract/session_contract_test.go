package contract

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestGetSessionContract_T025(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Session Contract User",
		DateOfBirth:     "1990-01-15",
		MobileNumber:    "+1987654329",
		Email:           "sessioncontract@example.com",
		Username:        "sessioncontract",
		Password:        "SessionPass123!",
		ConfirmPassword: "SessionPass123!",
	}
	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Login
	loginResp, err := provider.Login(auth.LoginRequest{Username: "sessioncontract", Password: "SessionPass123!"})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Validate getting the session
	sess, err := provider.GetSession(loginResp.SessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if sess.ID != loginResp.SessionID {
		t.Errorf("Expected Session ID %s, got %s", loginResp.SessionID, sess.ID)
	}
	if sess.UserID == "" {
		t.Error("Expected valid UserID in session")
	}
}
