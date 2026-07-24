package integration

import (
	"math"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestTimingAttackPrevention(t *testing.T) {
	provider := auth.NewLocalAuthProvider()
	// Register a valid user
	req := auth.RegisterUserRequest{
		FullName:        "Timing Attack User",
		DateOfBirth:     "1990-01-01",
		MobileNumber:    "9999999999",
		Email:           "timing@example.com",
		Username:        "timing_user",
		Password:        "Strong@Pass1234",
		ConfirmPassword: "Strong@Pass1234",
	}
	_, err := provider.Register(req)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Warmup
	provider.Login(auth.LoginRequest{Username: "timing_user", Password: "WrongPassword1!"})
	provider.Login(auth.LoginRequest{Username: "nonexistent_user", Password: "WrongPassword1!"})

	// Measure valid username but wrong password
	start := time.Now()
	for i := 0; i < 5; i++ {
		provider.Login(auth.LoginRequest{Username: "timing_user", Password: "WrongPassword1!"})
	}
	validUserDuration := time.Since(start) / 5

	// Measure invalid username and wrong password
	start2 := time.Now()
	for i := 0; i < 5; i++ {
		provider.Login(auth.LoginRequest{Username: "nonexistent_user", Password: "WrongPassword1!"})
	}
	invalidUserDuration := time.Since(start2) / 5

	diff := math.Abs(float64(validUserDuration - invalidUserDuration))
	// Difference should be small, say less than 20ms
	if diff > float64(20*time.Millisecond) {
		t.Errorf("Potential timing attack vulnerability: valid user took %v, invalid user took %v (diff: %v)", validUserDuration, invalidUserDuration, time.Duration(diff))
	}
}
