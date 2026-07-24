package integration

import (
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestSessionActiveExpiration_T032A(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	svc := sessions.NewService(repo, 2*time.Hour)
	provider := auth.NewLocalAuthProviderWithRepo(repo, 2*time.Hour)

	// Register user
	req := auth.RegisterUserRequest{
		FullName:        "Active Expiry User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999995",
		Email:           "expiry@example.com",
		Username:        "expiry_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	// Login
	resp, err := provider.Login(auth.LoginRequest{Username: "expiry_user", Password: "Strong@Pass1234"})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Get session directly from service
	sess, err := svc.GetSession(resp.SessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	originalActivity := sess.LastActivityAt

	// Simulate refreshing activity
	time.Sleep(1 * time.Millisecond) // Ensure time moves forward
	err = svc.RefreshActivity(sess.ID)
	if err != nil {
		t.Fatalf("Failed to refresh activity: %v", err)
	}

	updatedSess, err := svc.GetSession(sess.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if !updatedSess.LastActivityAt.After(originalActivity) {
		t.Errorf("Expected LastActivityAt to be updated. old: %v, new: %v", originalActivity, updatedSess.LastActivityAt)
	}
}
