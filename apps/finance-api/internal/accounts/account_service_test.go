package accounts

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/database"
	"artha-kosha/apps/finance-api/internal/domain"
)

type mockRepo struct {
	err error
	acc *Account
}

func (m *mockRepo) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.acc, nil
}

// mockAudit and mockDomain
type mockAudit struct {
	err error
}
func (m *mockAudit) RecordTx(ctx context.Context, tx *sql.Tx, ev audit.AuditEvent) error { return m.err }
func (m *mockAudit) Record(ctx context.Context, ev audit.AuditEvent) error { return m.err }

func TestCreateAccount_NoDB(t *testing.T) {
	repo := &mockRepo{acc: &Account{AccountID: "acc-1"}}
	svc := NewService(repo)

	acc, err := svc.CreateAccount(context.Background(), CreateAccountRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if acc.AccountID != "acc-1" {
		t.Errorf("wrong account ID")
	}
}

func TestCreateAccount_WithDB(t *testing.T) {
	repo := &mockRepo{acc: &Account{AccountID: "acc-2"}}
	
	// Create real audit/domain services but we won't execute them because tx will be nil or we mock database.WithTx
	// Wait, audit.Service uses SQL too, we can't just pass mockAudit easily since SetAuditService takes *audit.Service
	// But we can just pass nil to them or configure them with dummy dependencies if they don't panic.
	
	origWithTx := database.WithTx
	defer func() { database.WithTx = origWithTx }()

	// Test 1: WithTx succeeds, no audit, no domain
	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	
	svc := NewServiceWithDB(repo, &sql.DB{})
	
	acc, err := svc.CreateAccount(context.Background(), CreateAccountRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if acc.AccountID != "acc-2" {
		t.Errorf("wrong id")
	}

	// Test 2: Repo fails
	repo.err = errors.New("repo fail")
	_, err = svc.CreateAccount(context.Background(), CreateAccountRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	repo.err = nil

	// Test 3: WithTx fails
	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return errors.New("tx fail")
	}
	_, err = svc.CreateAccount(context.Background(), CreateAccountRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	
	// Test 4: WithTx succeeds, audit/domain set
	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	auditSvc := audit.NewService(&mockAuditRepo{})
	svc.SetAuditService(auditSvc)
	svc.SetDomainService(&mockDomainEmitter{})

	acc, err = svc.CreateAccount(context.Background(), CreateAccountRequest{OwnerID: "123"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if acc.AccountID != "acc-2" {
		t.Errorf("wrong id")
	}
}

type mockAuditRepo struct{}
func (m *mockAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }
func (m *mockAuditRepo) InsertTx(ctx context.Context, tx *sql.Tx, e audit.AuditEvent) error { return nil }

type mockDomainEmitter struct{}
func (m *mockDomainEmitter) EmitTx(ctx context.Context, tx *sql.Tx, e domain.DomainEvent) error { return nil }

