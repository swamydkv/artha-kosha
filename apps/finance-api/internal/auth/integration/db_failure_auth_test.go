package integration

import (
	"strings"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

// mockFailingRepo simulates a database connection failure for tests
type mockFailingRepo struct{}

func (m *mockFailingRepo) SaveSession(s sessions.Session) error {
	return sessions.ErrSessionNotFound // Just return some error to simulate failure
}
func (m *mockFailingRepo) GetSession(id string) (sessions.Session, error) {
	return sessions.Session{}, sessions.ErrSessionNotFound
}
func (m *mockFailingRepo) DeleteSession(id string) error {
	return nil
}
func (m *mockFailingRepo) DeleteAllUserSessions(userID string) error {
	return nil
}
func (m *mockFailingRepo) DeleteExpiredSessions() error {
	return nil
}

func TestDBFailureDuringAuthentication(t *testing.T) {
	// Create a provider with a failing session repository to simulate DB failure during login (session creation)
	provider := auth.NewLocalAuthProviderWithRepo(&mockFailingRepo{}, 1*time.Hour)

	// Since NewLocalAuthProviderWithRepo uses in-memory users, we can inject a user
	// by first using a normal provider to register, then migrating the user over, 
	// but the easiest way is just to register directly (register doesn't use session repo, it only uses it for DB in real postgres).
	// In the real implementation, NewLocalAuthProviderFromDSN is used for Postgres, but here we just mock the session repo failure.

	req := auth.RegisterUserRequest{
		FullName:        "DB Fail User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999998",
		Email:           "dbfail@example.com",
		Username:        "dbfail_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}

	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Now try to login. The password check will succeed, but session creation will fail.
	_, err = provider.Login(auth.LoginRequest{
		Username: "dbfail_user",
		Password: "Strong@Pass1234",
	})

	if err == nil {
		t.Fatal("Expected login to fail due to DB connection failure, but it succeeded")
	}
	
	// The error returned should ideally obscure internal DB details or just be an internal server error.
	if strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("Expected internal error for DB failure, got invalid credentials")
	}
}
