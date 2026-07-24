package integration

import (
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestSessionCreationAndValidation_T028(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	svc := sessions.NewService(repo, 2*time.Hour)
	provider := auth.NewLocalAuthProviderWithRepo(repo, 2*time.Hour)

	req := auth.RegisterUserRequest{
		FullName:        "Flow User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999996",
		Email:           "flow@example.com",
		Username:        "flow_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	resp, err := provider.Login(auth.LoginRequest{Username: "flow_user", Password: "Strong@Pass1234"})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if resp.SessionID == "" {
		t.Fatal("Expected non-empty session ID")
	}

	sess, err := svc.GetSession(resp.SessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if sess.UserID == "" || sess.Status != sessions.StatusActive {
		t.Errorf("Invalid session state: %+v", sess)
	}
}
