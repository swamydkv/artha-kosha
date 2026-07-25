package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/users"
)

type mockUsersRepo struct {
	err     error
	sessErr error
}

func (m mockUsersRepo) CreateUser(ctx context.Context, req users.CreateUserRequest) (*users.User, error) {
	if m.err != nil { return nil, m.err }
	return &users.User{UserID: "u1", Username: req.Username, PasswordHash: req.PasswordHash}, nil
}
func (m mockUsersRepo) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	if m.err != nil { return nil, m.err }
	hash, _ := hashPassword("Test1234!@#$")
	return &users.User{UserID: "u1", Username: username, FullName: "John Doe", PasswordHash: hash}, nil
}
func (m mockUsersRepo) GetUserByEmail(ctx context.Context, email string) (*users.User, error) { return nil, nil }
func (m mockUsersRepo) GetUserByMobileNumber(ctx context.Context, mobile string) (*users.User, error) { return nil, nil }
func (m mockUsersRepo) GetUserByID(ctx context.Context, id string) (*users.User, error) { return nil, nil }
func (m mockUsersRepo) CreateSession(ctx context.Context, sessionID, userID, ip, ua string) error { return m.sessErr }
func (m mockUsersRepo) DeleteUser(ctx context.Context, userID string, ret int) error { return m.err }
func (m mockUsersRepo) PruneArchivedUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m mockUsersRepo) CheckUserExists(ctx context.Context, email, username, mobile string) (*users.UserExistsCheck, error) { return &users.UserExistsCheck{}, nil }

func TestLocalAuthProvider_Register(t *testing.T) {
	p := NewLocalAuthProvider()
	// p.SetDomainService(nil)

	req := RegisterUserRequest{
		FullName:        "John Doe",
		DateOfBirth:     "2000-01-01",
		MobileNumber:    "1234567890",
		Email:           "john@example.com",
		Username:        "johndoe",
		Password:        "Test1234!@#$",
		ConfirmPassword: "Test1234!@#$",
	}

	// 1. In-memory
	res, err := p.Register(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if res.Username != "johndoe" { t.Errorf("expected johndoe, got %s", res.Username) }

	// Duplicate
	_, err = p.Register(req)
	if err == nil { t.Errorf("expected duplicate username error") }

	// 2. DB mode
	p.usersRepo = mockUsersRepo{}
	req.Username = "janedoe"
	res, err = p.Register(req)
	// DB error
	p.usersRepo = mockUsersRepo{err: errors.New("db err")}
	_, err = p.Register(req)
	if err == nil { t.Errorf("expected error for db error") }
}

func TestLocalAuthProvider_Login(t *testing.T) {
	p := NewLocalAuthProvider()
	// p.SetDomainService(nil)

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

	loginReq := LoginRequest{
		Username: "johndoe",
		Password: "Test1234!@#$",
	}

	// In memory success
	res, err := p.Login(loginReq)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if res.SessionID == "" { t.Errorf("expected session id") }

	// In memory invalid password
	loginReq.Password = "Wrong1234!@#$"
	_, err = p.Login(loginReq)
	if err == nil { t.Errorf("expected invalid credentials") }

	// DB mode
	p.usersRepo = mockUsersRepo{}
	loginReq.Password = "Test1234!@#$"
	res, err = p.Login(loginReq)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if res.SessionID == "" { t.Errorf("expected session id") }

	// DB get user error
	p.usersRepo = mockUsersRepo{err: errors.New("db error")}
	_, err = p.Login(loginReq)
	if err == nil { t.Errorf("expected error when get user fails") }

	// DB create session error
	p.usersRepo = mockUsersRepo{sessErr: errors.New("db session error")}
	_, err = p.Login(loginReq)
	if err == nil { t.Errorf("expected error when create session fails") }
}

func TestLocalAuthProvider_DeleteUser(t *testing.T) {
	p := NewLocalAuthProvider()
	
	// Error in memory
	err := p.DeleteUser("u1", "DELETE", 30)
	if err == nil { t.Errorf("expected error in memory mode") }

	// Error invalid confirmation
	p.usersRepo = mockUsersRepo{}
	err = p.DeleteUser("u1", "del", 30)
	if err == nil { t.Errorf("expected error for invalid confirmation") }

	// Success in DB
	err = p.DeleteUser("u1", "DELETE", 30)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}

func TestLocalAuthProvider_GetSession(t *testing.T) {
	p := NewLocalAuthProvider()
	req := RegisterUserRequest{
		FullName:        "John",
		DateOfBirth:     "2000-01-01",
		MobileNumber:    "1234567890",
		Email:           "john@example.com",
		Username:        "johndoe",
		Password:        "Test1234!@#$",
		ConfirmPassword: "Test1234!@#$",
	}
	_, err := p.Register(req)
	if err != nil { t.Fatalf("register failed: %v", err) }
	res, err := p.Login(LoginRequest{Username: "johndoe", Password: "Test1234!@#$"})
	if err != nil { t.Fatalf("login failed: %v", err) }

	sess, err := p.GetSession(res.SessionID)
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if sess.ID != res.SessionID { t.Errorf("expected session id") }

	// Logout
	_ = p.Logout(res.SessionID)
	_, err = p.GetSession(res.SessionID)
	if err == nil { t.Errorf("expected error on revoked session") }
}

func TestLocalAuthProvider_NewLocalAuthProviderFromDSN(t *testing.T) {
	// Should fail due to bad DSN
	_, _, err := NewLocalAuthProviderFromDSN("invalid", time.Hour)
	if err == nil { t.Errorf("expected error") }
}

func TestValidation(t *testing.T) {
	err := validateRegistrationRequest(RegisterUserRequest{})
	if err == nil { t.Errorf("expected error") }

	err = validateRegistrationRequest(RegisterUserRequest{
		FullName: "A", DateOfBirth: "A", MobileNumber: "A", Email: "A", Username: "A", Password: "A", ConfirmPassword: "B",
	})
	if err == nil { t.Errorf("expected error") }

	err = validateRegistrationRequest(RegisterUserRequest{
		FullName: "A", DateOfBirth: "A", MobileNumber: "A", Email: "A", Username: "A", Password: "A", ConfirmPassword: "A",
	})
	if err == nil { t.Errorf("expected error") }

	err = validateRegistrationRequest(RegisterUserRequest{
		FullName: "A", DateOfBirth: "2200-01-01", MobileNumber: "A", Email: "a@a.com", Username: "usrname", Password: "Password1234!", ConfirmPassword: "Password1234!",
	})
	if err == nil { t.Errorf("expected error") }
}
