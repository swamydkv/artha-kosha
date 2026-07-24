package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/users"
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

func TestProvider_Integration_RegistrationAndLogin(t *testing.T) {
	t.Log("Integration test for user registration and login flows using PostgreSQL")
	t.Skip("Skipping DB integration test in unit test file")
}

type mockUsersRepo struct {
	shouldErr bool
}

func (m *mockUsersRepo) CreateUser(ctx context.Context, req users.CreateUserRequest) (*users.User, error) {
	if m.shouldErr {
		return nil, errors.New("db error")
	}
	return &users.User{UserID: "user-1", Username: req.Username}, nil
}

func (m *mockUsersRepo) GetUserByID(ctx context.Context, id string) (*users.User, error) {
	return nil, nil
}
func (m *mockUsersRepo) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	if m.shouldErr {
		return nil, errors.New("not found")
	}
	hash, _ := hashPassword("Password123!")
	return &users.User{UserID: "user-1", Username: username, PasswordHash: hash}, nil
}
func (m *mockUsersRepo) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *mockUsersRepo) GetUserByMobileNumber(ctx context.Context, mobile string) (*users.User, error) {
	return nil, nil
}
func (m *mockUsersRepo) CheckUserExists(ctx context.Context, username, email, mobile string) (*users.UserExistsCheck, error) {
	return &users.UserExistsCheck{}, nil
}
func (m *mockUsersRepo) CreateSession(ctx context.Context, sessionID, userID, ipAddress, userAgent string) error {
	return nil
}

func TestProvider_DBMode_Branches(t *testing.T) {
	repo := &mockUsersRepo{}
	p := NewLocalAuthProvider()
	p.usersRepo = repo

	req := RegisterUserRequest{
		FullName: "Test DB",
		DateOfBirth: "1990-01-01",
		MobileNumber: "+1999999999",
		Email: "db@example.com",
		Username: "dbuser",
		Password: "Password123!",
		ConfirmPassword: "Password123!",
	}
	
	_, err := p.Register(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	repo.shouldErr = true
	_, err = p.Register(req)
	if err == nil {
		t.Error("expected error from db mock")
	}
	repo.shouldErr = false

	lReq := LoginRequest{
		Username: "dbuser",
		Password: "Password123!",
	}
	_, err = p.Login(lReq)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	repo.shouldErr = true
	_, err = p.Login(lReq)
	if err == nil {
		t.Error("expected error from db mock")
	}

	// validation errors
	badReq := req
	badReq.Password = "short"
	_, err = p.Register(badReq)
	if err == nil {
		t.Error("expected error for bad password")
	}

	badReq.Username = "!!"
	_, err = p.Register(badReq)
	if err == nil {
		t.Error("expected error for bad username")
	}
	
	badReq.DateOfBirth = "2999-01-01"
	_, err = p.Register(badReq)
	if err == nil {
		t.Error("expected error for future DOB")
	}
	
	badReq.DateOfBirth = "invalid"
	_, err = p.Register(badReq)
	if err == nil {
		t.Error("expected error for invalid DOB")
	}
	
	badReq.FullName = " "
	_, err = p.Register(badReq)
	if err == nil {
		t.Error("expected error for empty field")
	}
}

