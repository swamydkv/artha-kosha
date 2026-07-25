package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/domain"
	"artha-kosha/apps/finance-api/internal/outbox"
	"artha-kosha/apps/finance-api/internal/sessions"
	"artha-kosha/apps/finance-api/internal/transactions"
)

type mockAuthProvider struct {
	err error
	acc *accounts.Service
	tx  *transactions.Service
	bdg *budgets.Service
}

func (m *mockAuthProvider) Register(req RegisterUserRequest) (RegisterUserResponse, error) {
	return RegisterUserResponse{}, m.err
}
func (m *mockAuthProvider) Login(req LoginRequest) (LoginResponse, error) {
	if m.err != nil {
		return LoginResponse{}, m.err
	}
	return LoginResponse{
		UserID:    uuid.New().String(),
		SessionID: uuid.New().String(),
	}, nil
}
func (m *mockAuthProvider) Logout(sessionID string) error { return m.err }
func (m *mockAuthProvider) GetSession(sessionID string) (sessions.Session, error) {
	return sessions.Session{}, m.err
}
func (m *mockAuthProvider) RevokeAll(userID string) error { return m.err }
func (m *mockAuthProvider) DeleteUser(userID string, confirmation string, archiveRetentionDays int) error {
	return m.err
}
func (m *mockAuthProvider) ChangePassword(req ChangePasswordRequest) error { return m.err }

func (m *mockAuthProvider) GetAccountsService() *accounts.Service         { return m.acc }
func (m *mockAuthProvider) GetTransactionsService() *transactions.Service { return m.tx }
func (m *mockAuthProvider) GetBudgetsService() *budgets.Service           { return m.bdg }

func TestLoginService(t *testing.T) {
	provider := &mockAuthProvider{}
	svc := NewLoginService(provider, nil)

	_, _ = svc.Login(LoginRequest{Username: "a", Password: "b"})
	_ = svc.Logout("s1")
	_, _ = svc.Register(RegisterUserRequest{})
	_, _ = svc.GetSession("s1")
	_ = svc.RevokeAll("u1")
	_ = svc.ChangePassword(ChangePasswordRequest{})
	_ = svc.DeleteUser("u1", "DELETE", 30)

	_ = svc.GetAccountsService()
	_ = svc.GetTransactionsService()
	_ = svc.GetBudgetsService()
}

func TestLoginServiceWithDB_LoginCtx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	provider := &mockAuthProvider{}
	svc := NewLoginServiceWithDB(provider, nil, db)

	// Mock domain and audit services
	auditRepo := audit.NewSQLRepository(db)
	auditSvc := audit.NewService(auditRepo)
	svc.SetAuditService(auditSvc)

	domainRepo := domain.NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	domainSvc := domain.NewService(domainRepo, outboxSvc)
	svc.SetDomainService(domainSvc)

	t.Run("empty credentials", func(t *testing.T) {
		_, err := svc.LoginCtx(context.Background(), LoginRequest{})
		if err == nil {
			t.Error("expected error for empty credentials")
		}
	})

	t.Run("successful login with tx", func(t *testing.T) {
		mock.ExpectBegin()
		// Audit event
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		// Domain event
		mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
		// Outbox entry
		mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := svc.LoginCtx(context.Background(), LoginRequest{Username: "a", Password: "b"})
		if err != nil {
			t.Errorf("expected nil err, got %v", err)
		}
	})

	t.Run("failed provider login", func(t *testing.T) {
		provider.err = errors.New("auth failed")
		mock.ExpectBegin()
		mock.ExpectRollback()

		_, err := svc.LoginCtx(context.Background(), LoginRequest{Username: "a", Password: "b"})
		if err == nil {
			t.Error("expected error")
		}
		provider.err = nil
	})

	t.Run("audit failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		_, err := svc.LoginCtx(context.Background(), LoginRequest{Username: "a", Password: "b"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("domain event failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO domain_events").WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		_, err := svc.LoginCtx(context.Background(), LoginRequest{Username: "a", Password: "b"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("commit failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		_, err := svc.LoginCtx(context.Background(), LoginRequest{Username: "a", Password: "b"})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestLoginServiceWithDB_LogoutCtx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	provider := &mockAuthProvider{}
	svc := NewLoginServiceWithDB(provider, nil, db)

	// Mock domain and audit services
	auditRepo := audit.NewSQLRepository(db)
	auditSvc := audit.NewService(auditRepo)
	svc.SetAuditService(auditSvc)

	domainRepo := domain.NewSQLRepository(db)
	outboxRepo := outbox.NewSQLRepository(db)
	outboxSvc := outbox.NewService(outboxRepo)
	domainSvc := domain.NewService(domainRepo, outboxSvc)
	svc.SetDomainService(domainSvc)

	t.Run("empty session ID", func(t *testing.T) {
		err := svc.LogoutCtx(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty session id")
		}
	})

	t.Run("successful logout with tx", func(t *testing.T) {
		mock.ExpectBegin()
		// Audit event
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		// Domain event
		mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
		// Outbox entry
		mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := svc.LogoutCtx(context.Background(), "s1")
		if err != nil {
			t.Errorf("expected nil err, got %v", err)
		}
	})

	t.Run("failed provider logout", func(t *testing.T) {
		provider.err = errors.New("auth failed")
		mock.ExpectBegin()
		mock.ExpectRollback()

		err := svc.LogoutCtx(context.Background(), "s1")
		if err == nil {
			t.Error("expected error")
		}
		provider.err = nil
	})

	t.Run("audit failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := svc.LogoutCtx(context.Background(), "s1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("domain event failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO domain_events").WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := svc.LogoutCtx(context.Background(), "s1")
		if err == nil {
			t.Error("expected error")
		}
	})
	
	t.Run("commit failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO domain_events").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO transactional_outbox").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

		err := svc.LogoutCtx(context.Background(), "s1")
		if err == nil {
			t.Error("expected error")
		}
	})
	
	t.Run("no db fallback login/logout", func(t *testing.T) {
		svcNoDB := NewLoginService(provider, nil)
		svcNoDB.SetAuditService(auditSvc)
		mock.ExpectExec("INSERT INTO audit_events").WillReturnResult(sqlmock.NewResult(1, 1))
		_ = svcNoDB.Logout("s1")
		_, _ = svcNoDB.Login(LoginRequest{Username: "a", Password: "b"})
	})
	
	t.Run("no db fallback provider err", func(t *testing.T) {
		svcNoDB := NewLoginService(provider, nil)
		provider.err = errors.New("err")
		_, _ = svcNoDB.Login(LoginRequest{Username: "a", Password: "b"})
		_ = svcNoDB.Logout("s1")
		provider.err = nil
	})
}

type mockAuthProviderNoMethods struct {
	err error
}
func (m *mockAuthProviderNoMethods) Register(req RegisterUserRequest) (RegisterUserResponse, error) { return RegisterUserResponse{}, m.err }
func (m *mockAuthProviderNoMethods) Login(req LoginRequest) (LoginResponse, error) { return LoginResponse{}, m.err }
func (m *mockAuthProviderNoMethods) Logout(sessionID string) error { return m.err }
func (m *mockAuthProviderNoMethods) GetSession(sessionID string) (sessions.Session, error) { return sessions.Session{}, m.err }
func (m *mockAuthProviderNoMethods) RevokeAll(userID string) error { return m.err }
func (m *mockAuthProviderNoMethods) DeleteUser(userID string, confirmation string, archiveRetentionDays int) error { return m.err }
func (m *mockAuthProviderNoMethods) ChangePassword(req ChangePasswordRequest) error { return m.err }

func TestLoginService_NoMethodsFallback(t *testing.T) {
	provider := &mockAuthProviderNoMethods{}
	svc := NewLoginService(provider, nil)

	if svc.GetAccountsService() != nil {
		t.Error("expected nil")
	}
	if svc.GetTransactionsService() != nil {
		t.Error("expected nil")
	}
	if svc.GetBudgetsService() != nil {
		t.Error("expected nil")
	}
}
