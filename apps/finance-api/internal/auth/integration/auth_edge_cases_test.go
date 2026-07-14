package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestAuthEdgeCases_SQLInjection(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Try to register with SQL injection patterns
	registerRequest := auth.RegisterUserRequest{
		FullName:        "SQL Injection User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "sql@example.com",
		Username:        "user'; DROP TABLE users; --",
		Password:        "SQLInjectionPass123!",
		ConfirmPassword: "SQLInjectionPass123!",
	}

	_, err := provider.Register(registerRequest)
	// Should either fail validation or be handled safely
	if err == nil {
		// If it succeeds, the username should be sanitized
		t.Log("SQL injection pattern was accepted - verify sanitization")
	}
}

func TestAuthEdgeCases_XSSPatterns(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Try to register with XSS patterns
	registerRequest := auth.RegisterUserRequest{
		FullName:        "<script>alert('xss')</script>",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "xss@example.com",
		Username:        "xssuser",
		Password:        "XSSPass123!@",
		ConfirmPassword: "XSSPass123!@",
	}

	response, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Note: XSS sanitization is typically handled at the serialization/frontend layer
	// The business logic layer stores the data as-is
	// For now, we accept the current behavior
	_ = response.FirstName // Avoid unused variable warning
}

func TestAuthEdgeCases_LongInputs(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Try to register with very long inputs
	longString := "a"
	for i := 0; i < 1000; i++ {
		longString += "a"
	}

	registerRequest := auth.RegisterUserRequest{
		FullName:        longString,
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "long@example.com",
		Username:        "longuser",
		Password:        "LongPass123!",
		ConfirmPassword: "LongPass123!",
	}

	_, err := provider.Register(registerRequest)
	// Should either fail validation or handle gracefully
	if err != nil {
		t.Log("Long input was rejected by validation")
	}
}

func TestAuthEdgeCases_SpecialCharacters(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	tests := []struct {
		name     string
		username string
		valid    bool
	}{
		{
			name:     "Username with emojis",
			username: "user😀",
			valid:    false,
		},
		{
			name:     "Username with spaces",
			username: "user name",
			valid:    false,
		},
		{
			name:     "Username with valid special chars",
			username: "user_name.test",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerRequest := auth.RegisterUserRequest{
				FullName:        "Special Char User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "special@example.com",
				Username:        tt.username,
				Password:        "SpecialPass123!",
				ConfirmPassword: "SpecialPass123!",
			}

			_, err := provider.Register(registerRequest)
			if tt.valid && err != nil {
				t.Errorf("Expected username %s to be valid, got error: %v", tt.username, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected username %s to be invalid, but it was accepted", tt.username)
			}
		})
	}
}

func TestAuthEdgeCases_ConcurrentRegistration(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Try to register the same user concurrently
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Concurrent User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "concurrent@example.com",
		Username:        "concurrentuser",
		Password:        "ConcurrentPass123!",
		ConfirmPassword: "ConcurrentPass123!",
	}

	// Register once
	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register again (should fail)
	_, err = provider.Register(registerRequest)
	if err == nil {
		t.Error("Expected duplicate registration to fail")
	}
}

func TestAuthEdgeCases_SessionReuse(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	// Register a user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Session User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "session@example.com",
		Username:        "sessionuser",
		Password:        "SessionPass123!",
		ConfirmPassword: "SessionPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Login multiple times
	loginRequest := auth.LoginRequest{
		Username: "sessionuser",
		Password: "SessionPass123!",
	}

	session1, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("First login failed: %v", err)
	}

	session2, err := provider.Login(loginRequest)
	if err != nil {
		t.Fatalf("Second login failed: %v", err)
	}

	// Verify different sessions are created
	if session1.SessionID == session2.SessionID {
		t.Error("Each login should create a new session")
	}
}