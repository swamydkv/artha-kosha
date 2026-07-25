package auth

import (
	"testing"
)

func TestLocalAuthProvider_ChangePassword(t *testing.T) {
	p := NewLocalAuthProvider()

	// Setup a user
	req := RegisterUserRequest{
		FullName:        "John Doe",
		DateOfBirth:     "2000-01-01",
		MobileNumber:    "1234567890",
		Email:           "john@example.com",
		Username:        "johndoe",
		Password:        "Test1234!@#$",
		ConfirmPassword: "Test1234!@#$",
	}
	_, _ = p.Register(req)

	// User not found
	err := p.ChangePassword(ChangePasswordRequest{Username: "unknown", OldPassword: "a", NewPassword: "b"})
	if err == nil { t.Fatalf("expected error for user not found") }

	// Invalid old password
	err = p.ChangePassword(ChangePasswordRequest{Username: "johndoe", OldPassword: "Wrong1234!@#$", NewPassword: "New1234!@#$"})
	if err == nil { t.Fatalf("expected error for invalid old password") }

	// Short new password
	err = p.ChangePassword(ChangePasswordRequest{Username: "johndoe", OldPassword: "Test1234!@#$", NewPassword: "short"})
	if err == nil { t.Fatalf("expected error for short new password") }

	// Success
	err = p.ChangePassword(ChangePasswordRequest{Username: "johndoe", OldPassword: "Test1234!@#$", NewPassword: "NewPass1234!@#$"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}
