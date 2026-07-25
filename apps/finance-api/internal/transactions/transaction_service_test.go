package transactions

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
	tx  *Transaction
}

func (m *mockRepo) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tx, nil
}

func TestCreateTransaction_NoDB(t *testing.T) {
	repo := &mockRepo{tx: &Transaction{TransactionID: "txn-1"}}
	svc := NewService(repo)

	tx, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if tx.TransactionID != "txn-1" {
		t.Errorf("wrong transaction ID")
	}
}

func TestCreateTransaction_WithDB(t *testing.T) {
	repo := &mockRepo{tx: &Transaction{TransactionID: "txn-2"}}
	
	origWithTx := database.WithTx
	defer func() { database.WithTx = origWithTx }()

	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	
	svc := NewServiceWithDB(repo, &sql.DB{})
	
	tx, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if tx.TransactionID != "txn-2" {
		t.Errorf("wrong id")
	}

	repo.err = errors.New("repo fail")
	_, err = svc.CreateTransaction(context.Background(), CreateTransactionRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	repo.err = nil

	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return errors.New("tx fail")
	}
	_, err = svc.CreateTransaction(context.Background(), CreateTransactionRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}

	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	// Test 4: WithTx succeeds, audit/domain set
	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	auditSvc := audit.NewService(&mockAuditRepo{})
	svc.SetAuditService(auditSvc)
	svc.SetDomainService(&mockDomainEmitter{})

	txObj, err := svc.CreateTransaction(context.Background(), CreateTransactionRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if txObj.TransactionID != "txn-2" {
		t.Errorf("wrong id")
	}
}

type mockAuditRepo struct{}
func (m *mockAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }
func (m *mockAuditRepo) InsertTx(ctx context.Context, tx *sql.Tx, e audit.AuditEvent) error { return nil }

type mockDomainEmitter struct{}
func (m *mockDomainEmitter) EmitTx(ctx context.Context, tx *sql.Tx, e domain.DomainEvent) error { return nil }
