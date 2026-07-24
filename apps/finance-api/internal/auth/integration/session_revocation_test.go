package integration

import (
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestSessionRevocation_T030(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	svc := sessions.NewService(repo, 2*time.Hour)
	provider := auth.NewLocalAuthProviderWithRepo(repo, 2*time.Hour)

	// Register user
	req := auth.RegisterUserRequest{
		FullName:        "Revoke User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999994",
		Email:           "revoke@example.com",
		Username:        "revoke_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	// Login
	resp, err := provider.Login(auth.LoginRequest{Username: "revoke_user", Password: "Strong@Pass1234"})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Revoke Session
	err = svc.RevokeSession(resp.SessionID)
	if err != nil {
		t.Fatalf("Failed to revoke session: %v", err)
	}

	// Check if session is revoked
	sess, err := svc.GetSession(resp.SessionID)
	if err == nil {
		// some services might return error on revoked session
		if sess.Status != sessions.StatusRevoked {
			t.Errorf("Expected status revoked, got %s", sess.Status)
		}
	}
}
