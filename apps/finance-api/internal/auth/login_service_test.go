package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"artha-kosha/apps/finance-api/internal/sqlc"
	"artha-kosha/apps/finance-api/internal/sessions"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoginService_Login(t *testing.T) {
	svc := NewLoginService(&mockProvider{}, nil)
	
	req := LoginRequest{}
	_, err := svc.Login(context.Background(), req)
	if err == nil {
		t.Error("expected error on login")
	}

	req.Username = "user"
	req.Password = "pass"
	_, err = svc.Login(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// test db
	db, mock, _ := sqlmock.New()
	svcDB := NewLoginServiceWithDB(&mockProvider{}, nil, db)
	svcDB.SetAuditService(nil)
	svcDB.SetDomainService(nil)
	
	mock.ExpectBegin()
	mock.ExpectCommit()
	_, err = svcDB.Login(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error with DB, got %v", err)
	}
}

func TestLoginService_Logout(t *testing.T) {
	svc := NewLoginService(&mockProvider{}, nil)
	err := svc.Logout(context.Background(), "session")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	db, mock, _ := sqlmock.New()
	svcDB := NewLoginServiceWithDB(&mockProvider{}, nil, db)
	mock.ExpectBegin()
	mock.ExpectCommit()
	err = svcDB.Logout(context.Background(), "session")
	if err != nil {
		t.Errorf("expected no error with DB, got %v", err)
	}
}

func TestLoginService_Logout_EmptySession(t *testing.T) {
	svc := NewLoginService(&mockProvider{}, nil)
	err := svc.Logout(context.Background(), "")
	if err == nil {
		t.Error("expected error on empty session")
	}
}

func TestLoginService_genIDLogin(t *testing.T) {
	id := genIDLogin("test")
	if id == "" {
		t.Error("expected non-empty id")
	}
}

// Mock provider
type mockProvider struct {}

func (m *mockProvider) Authenticate(ctx context.Context, username string, password string) (sqlc.User, error) {
	return sqlc.User{UserID: uuid.New()}, nil
}
func (m *mockProvider) CreateSession(ctx context.Context, userID uuid.UUID, clientIP string, userAgent string) (sqlc.Session, error) {
	return sqlc.Session{ID: uuid.New()}, nil
}
func (m *mockProvider) ValidateSession(ctx context.Context, sessionToken string) (sqlc.Session, error) {
	return sqlc.Session{}, nil
}
func (m *mockProvider) RevokeSession(ctx context.Context, sessionToken string) error {
	return nil
}
func (m *mockProvider) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (m *mockProvider) RefreshSession(ctx context.Context, sessionToken string) (sqlc.Session, error) {
	return sqlc.Session{}, nil
}
func (m *mockProvider) Register(req RegisterUserRequest) (RegisterUserResponse, error) {
	return RegisterUserResponse{}, nil
}
func (m *mockProvider) Login(req LoginRequest) (LoginResponse, error) {
	return LoginResponse{}, nil
}
func (m *mockProvider) Logout(sess string) error {
	return nil
}
func (m *mockProvider) GetSession(sess string) (sessions.Session, error) {
	return sessions.Session{}, nil
}
func (m *mockProvider) RevokeAll(id string) error {
	return nil
}
func (m *mockProvider) ChangePassword(req ChangePasswordRequest) error {
	return nil
}
