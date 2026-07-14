package auth

import (
	"strings"
	"testing"
)

func TestLoginReturnsGenericErrorForInvalidCredentials(t *testing.T) {
	provider := NewLocalAuthProvider()

	_, err := provider.Register(RegisterUserRequest{
		FullName:        "Jane Doe",
		DateOfBirth:     "1990-05-01",
		MobileNumber:    "+12065550123",
		Email:           "jane@example.com",
		Username:        "janedoe",
		Password:        "Passw0rd!123",
		ConfirmPassword: "Passw0rd!123",
	})
	if err != nil {
		t.Fatalf("expected initial registration to succeed, got %v", err)
	}

	_, err = provider.Login(LoginRequest{Username: "janedoe", Password: "wrong-password"})
	if err == nil {
		t.Fatal("expected invalid credentials to return an error")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("expected generic invalid credentials error, got %q", err.Error())
	}
}

func TestLoginCreatesSessionForValidCredentials(t *testing.T) {
	provider := NewLocalAuthProvider()

	_, err := provider.Register(RegisterUserRequest{
		FullName:        "Jane Doe",
		DateOfBirth:     "1990-05-01",
		MobileNumber:    "+12065550123",
		Email:           "jane@example.com",
		Username:        "janedoe",
		Password:        "Passw0rd!123",
		ConfirmPassword: "Passw0rd!123",
	})
	if err != nil {
		t.Fatalf("expected initial registration to succeed, got %v", err)
	}

	result, err := provider.Login(LoginRequest{Username: "janedoe", Password: "Passw0rd!123"})
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if result.SessionID == "" {
		t.Fatal("expected a session identifier to be generated")
	}
	if result.WelcomeMessage == "" {
		t.Fatal("expected a personalized welcome message")
	}
}
