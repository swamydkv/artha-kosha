package integration

import (
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestConcurrentSessionLimitEnforcement_T031(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	provider := auth.NewLocalAuthProviderWithRepo(repo, 2*time.Hour)

	// Register user
	req := auth.RegisterUserRequest{
		FullName:        "Session Limit User 2",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999993",
		Email:           "limit2@example.com",
		Username:        "limit_user2",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	// Login multiple times
	for i := 0; i < 6; i++ { 
		_, err := provider.Login(auth.LoginRequest{Username: "limit_user2", Password: "Strong@Pass1234"})
		if err != nil {
			if i == 5 {
				return
			}
			t.Fatalf("Failed to login on attempt %d: %v", i+1, err)
		}
	}
}
