package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func BenchmarkLogin(b *testing.B) {
	provider := auth.NewLocalAuthProvider()

	// Register a test user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Benchmark Login User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1888777666",
		Email:           "benchmarklogin@example.com",
		Username:        "benchmarkloginuser",
		Password:        "BenchmarkLoginPass123!",
		ConfirmPassword: "BenchmarkLoginPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		b.Fatalf("Failed to register test user: %v", err)
	}

	loginRequest := auth.LoginRequest{
		Username: "benchmarkloginuser",
		Password: "BenchmarkLoginPass123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = provider.Login(loginRequest)
	}
}

func BenchmarkLogout(b *testing.B) {
	provider := auth.NewLocalAuthProvider()

	// Register a test user
	registerRequest := auth.RegisterUserRequest{
		FullName:        "Benchmark Logout User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1777666555",
		Email:           "benchmarklogout@example.com",
		Username:        "benchmarklogoutuser",
		Password:        "BenchmarkLogoutPass123!",
		ConfirmPassword: "BenchmarkLogoutPass123!",
	}

	_, err := provider.Register(registerRequest)
	if err != nil {
		b.Fatalf("Failed to register test user: %v", err)
	}

	loginRequest := auth.LoginRequest{
		Username: "benchmarklogoutuser",
		Password: "BenchmarkLogoutPass123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session, _ := provider.Login(loginRequest)
		_ = provider.Logout(session.SessionID)
	}
}