package auth

import "testing"

func TestRegisterUserCreatesAccountAndHashesPassword(t *testing.T) {
	provider := NewLocalAuthProvider()

	result, err := provider.Register(RegisterUserRequest{
		FullName:        "Jane Doe",
		DateOfBirth:     "1990-05-01",
		MobileNumber:    "+12065550123",
		Email:           "jane@example.com",
		Username:        "janedoe",
		Password:        "Passw0rd!123",
		ConfirmPassword: "Passw0rd!123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.UserID == "" {
		t.Fatal("expected a generated user id")
	}
	if result.Username != "janedoe" {
		t.Fatalf("expected username janedoe, got %q", result.Username)
	}
	if result.PasswordHash == "" {
		t.Fatal("expected a password hash to be generated")
	}
	if result.PasswordHash == "Passw0rd!123" {
		t.Fatal("password hash must not equal the plaintext password")
	}
}

func TestRegisterUserRejectsDuplicateCredentials(t *testing.T) {
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
		t.Fatalf("expected first registration to succeed, got %v", err)
	}

	_, err = provider.Register(RegisterUserRequest{
		FullName:        "Jane Doe",
		DateOfBirth:     "1990-05-01",
		MobileNumber:    "+12065550123",
		Email:           "jane@example.com",
		Username:        "janedoe",
		Password:        "Passw0rd!123",
		ConfirmPassword: "Passw0rd!123",
	})
	if err == nil {
		t.Fatal("expected duplicate registration to fail")
	}
}
