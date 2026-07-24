package integration

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func BenchmarkRegistration(b *testing.B) {
	provider := auth.NewLocalAuthProvider()

	request := auth.RegisterUserRequest{
		FullName:        "Benchmark User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "benchmark@example.com",
		Username:        "benchmarkuser",
		Password:        "BenchmarkPass123!",
		ConfirmPassword: "BenchmarkPass123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use unique username/email for each iteration
		request.Username = "benchmarkuser" + string(rune(i))
		request.Email = "benchmark" + string(rune(i)) + "@example.com"
		_, _ = provider.Register(request)
	}
}

func BenchmarkRegistrationValidation(b *testing.B) {
	provider := auth.NewLocalAuthProvider()

	request := auth.RegisterUserRequest{
		FullName:        "Benchmark User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "+1234567890",
		Email:           "benchmark@example.com",
		Username:        "benchmarkuser",
		Password:        "BenchmarkPass123!",
		ConfirmPassword: "BenchmarkPass123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use unique username/email for each iteration
		request.Username = "benchmarkuser" + string(rune(i))
		request.Email = "benchmark" + string(rune(i)) + "@example.com"
		_, _ = provider.Register(request)
	}
}
