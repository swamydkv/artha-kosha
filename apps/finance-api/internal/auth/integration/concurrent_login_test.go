package integration

import (
	"sync"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestConcurrentLoginAttempts_T032C(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	provider := auth.NewLocalAuthProviderWithRepo(repo, 2*time.Hour)

	// Register user
	req := auth.RegisterUserRequest{
		FullName:        "Concurrent User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999991",
		Email:           "concurrent@example.com",
		Username:        "concurrent_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	var wg sync.WaitGroup
	errorsCh := make(chan error, 10)
	
	// Simulate 10 concurrent logins
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := provider.Login(auth.LoginRequest{Username: "concurrent_user", Password: "Strong@Pass1234"})
			if err != nil {
				errorsCh <- err
			}
		}()
	}
	wg.Wait()
	close(errorsCh)

	var errs []error
	for err := range errorsCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		t.Errorf("Expected all concurrent logins to succeed, but got %d errors. First error: %v", len(errs), errs[0])
	}
}
