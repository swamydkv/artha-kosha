package auth

import (
	"context"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"artha-kosha/apps/finance-api/internal/users"
)

func TestRegisterService(t *testing.T) {
	svc := NewRegisterService(&mockProvider{}, &mockUserRepo{})
	svc.SetDomainService(nil)
	svc.SetAuditService(nil)

	req := RegisterUserRequest{Email: "invalid"}
	_, err := svc.Register(context.Background(), req)
	if err == nil {
		t.Error("expected error on invalid request")
	}

	req = RegisterUserRequest{
		FullName: "Test User",
		DateOfBirth: "1990-01-01",
		MobileNumber: "+1234567890",
		Email: "test@example.com",
		Username: "testuser",
		Password: "Password123!",
        ConfirmPassword: "Password123!",
	}
	_, err = svc.Register(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	id := genIDRegister("test")
	if id == "" {
		t.Error("expected non-empty id")
	}
}

func TestRegisterServiceWithDB(t *testing.T) {
    db, mock, _ := sqlmock.New()
    svc := NewRegisterServiceWithDB(&mockProvider{}, &mockUserRepo{}, db)
    if svc == nil {
        t.Error("expected svc to be created")
    }
	
	req := RegisterUserRequest{
		FullName: "Test User",
		DateOfBirth: "1990-01-01",
		MobileNumber: "+1234567890",
		Email: "test@example.com",
		Username: "testuser",
		Password: "Password123!",
        ConfirmPassword: "Password123!",
	}
	mock.ExpectBegin()
	mock.ExpectCommit()
	_, err := svc.Register(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error with DB, got %v", err)
	}
}

type mockUserRepo struct{}
func (m *mockUserRepo) CheckUserExists(ctx context.Context, username string, email string, mobileNumber string) (*users.UserExistsCheck, error) {
	return &users.UserExistsCheck{}, nil
}
func (m *mockUserRepo) CreateUser(ctx context.Context, req users.CreateUserRequest) (*users.User, error) {
	return &users.User{UserID: "some-id"}, nil
}
func (m *mockUserRepo) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByMobileNumber(ctx context.Context, mobile string) (*users.User, error) {
	return nil, nil
}
