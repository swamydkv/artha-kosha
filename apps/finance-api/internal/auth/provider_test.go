package auth

import (
	"testing"
	"time"
)

func TestProviderSetters(t *testing.T) {
	p := &LocalAuthProvider{}
	p.SetDomainService(nil)
	p.SetAuditService(nil)
	p.SetAccountsService(nil)
	p.SetTransactionsService(nil)
	p.SetBudgetsService(nil)
}

func TestProvider_NewLocalAuthProviderFromDSN(t *testing.T) {
	_, _, err := NewLocalAuthProviderFromDSN("invalid-dsn", time.Minute)
	if err == nil {
		t.Error("expected error with invalid dsn")
	}
}

func TestProvider_firstNameFromFullName(t *testing.T) {
	if firstNameFromFullName("John Doe") != "John" {
		t.Error("expected John")
	}
	if firstNameFromFullName("Jane") != "Jane" {
		t.Error("expected Jane")
	}
}

func TestProvider_validateRegistrationRequest(t *testing.T) {
	req := RegisterUserRequest{
		FullName: "Test User",
		DateOfBirth: "1990-01-01",
		MobileNumber: "+1234567890",
		Email: "test@example.com",
		Username: "testuser",
		Password: "Password123!",
        ConfirmPassword: "Password123!",
	}
	
	err := validateRegistrationRequest(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	req.Email = "invalid"
	err = validateRegistrationRequest(req)
	if err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestProvider_AllBranches(t *testing.T) {
	p := NewLocalAuthProvider()
	// Test Register
	req := RegisterUserRequest{
		FullName: "Test User",
		DateOfBirth: "1990-01-01",
		MobileNumber: "+1234567890",
		Email: "test@example.com",
		Username: "testuser",
		Password: "Password123!",
        ConfirmPassword: "Password123!",
	}
	res, err := p.Register(req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	// duplicate register
	_, err = p.Register(req)
	if err == nil {
		t.Error("expected error on duplicate register")
	}
	// Test Login
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "Password123!",
	}
	lRes, err := p.Login(loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	// invalid login
	_, err = p.Login(LoginRequest{Username: "testuser", Password: "wrong"})
	if err == nil {
		t.Error("expected error on wrong password")
	}
	_, err = p.Login(LoginRequest{Username: "unknown", Password: "wrong"})
	if err == nil {
		t.Error("expected error on unknown user")
	}
	// Test ChangePassword
	err = p.ChangePassword(ChangePasswordRequest{
		Username: "testuser",
		OldPassword: "Password123!",
		NewPassword: "NewPassword123!",
	})
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}
	err = p.ChangePassword(ChangePasswordRequest{Username: "unknown", OldPassword: "", NewPassword: ""})
	if err == nil {
		t.Error("expected error on unknown user")
	}
	err = p.ChangePassword(ChangePasswordRequest{Username: "testuser", OldPassword: "wrong", NewPassword: ""})
	if err == nil {
		t.Error("expected error on wrong old password")
	}
	err = p.ChangePassword(ChangePasswordRequest{Username: "testuser", OldPassword: "NewPassword123!", NewPassword: "short"})
	if err == nil {
		t.Error("expected error on short new password")
	}
	
	// Logout
	err = p.Logout(lRes.SessionID)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
	// GetSession
	_, err = p.GetSession(lRes.SessionID)
	if err == nil {
		t.Error("expected error for revoked session")
	}
	// RevokeAll
	err = p.RevokeAll(res.UserID)
	if err != nil {
		t.Fatalf("RevokeAll failed: %v", err)
	}

	// extra validation branches
	req.Password = "noupper123!"
	req.ConfirmPassword = req.Password
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }
	req.Password = "NOLOWER123!"
	req.ConfirmPassword = req.Password
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }
	req.Password = "NoDigit!!!!"
	req.ConfirmPassword = req.Password
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }
	req.Password = "NoSpecial1234"
	req.ConfirmPassword = req.Password
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }
	req.Password = "Password123!"
	req.ConfirmPassword = req.Password
	req.Email = "nodomain@"
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }
	req.Email = "@nodomain.com"
	if validateRegistrationRequest(req) == nil { t.Error("expected error") }

	// getters
	p.GetAccountsService()
	p.GetTransactionsService()
	p.GetBudgetsService()

	// NewLocalAuthProviderWithRepo
	NewLocalAuthProviderWithRepo(nil, time.Minute)
}
