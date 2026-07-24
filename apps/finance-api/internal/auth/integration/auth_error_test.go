package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestAuthError_InvalidCredentials(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Error Test User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1222333444",
		Email:           "errortest@example.com",
		Username:        "errortestuser",
		Password:        "ErrorPass123!",
		ConfirmPassword: "ErrorPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "Wrong password",
			username: "errortestuser",
			password: "WrongPassword!",
		},
		{
			name:     "Non-existent user",
			username: "nonexistentuser",
			password: "SomePassword123!",
		},
		{
			name:     "Empty username",
			username: "",
			password: "ErrorPass123!",
		},
		{
			name:     "Empty password",
			username: "errortestuser",
			password: "",
		},
		{
			name:     "Both empty",
			username: "",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginRequest := auth.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			_, err := provider.Login(loginRequest)
			if err == nil {
				t.Error("Expected error for invalid credentials")
			}
		})
	}
}

func TestAuthError_LogoutInvalidSession(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Try to logout with invalid session ID
	err := provider.Logout("invalid-session-id")
	if err != nil {
		t.Errorf("Logout with invalid session should not error, got: %v", err)
	}
}

func TestAuthError_LoginAfterLogout(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        " Logout Test User",
		DateOfBirth:     "1995-06-15",
		MobileNumber:    "+1555666777",
		Email:           "logouttest@example.com",
		Username:        "logouttestuser",
		Password:        "LogoutPass123!",
		ConfirmPassword: "LogoutPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Login
	loginRequest := auth.LoginRequest{
		Username: "logouttestuser",
		Password: "LogoutPass123!",
	}

	session, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Logout
	err = provider.Logout(session.SessionID)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Login again should succeed
	newSession, err := provider.Login(loginRequest)
	if err != nil {
		t.Errorf("Login after logout should succeed, got error: %v", err)
	}

	if newSession.SessionID == session.SessionID {
		t.Error("New login should create a different session")
	}
}
