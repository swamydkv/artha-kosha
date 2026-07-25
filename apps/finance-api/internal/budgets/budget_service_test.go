package budgets

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
	b   *Budget
}

func (m *mockRepo) CreateBudget(ctx context.Context, req CreateBudgetRequest) (*Budget, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.b, nil
}

func TestCreateBudget_NoDB(t *testing.T) {
	repo := &mockRepo{b: &Budget{BudgetID: "bud-1"}}
	svc := NewService(repo)

	b, err := svc.CreateBudget(context.Background(), CreateBudgetRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if b.BudgetID != "bud-1" {
		t.Errorf("wrong budget ID")
	}
}

func TestCreateBudget_WithDB(t *testing.T) {
	repo := &mockRepo{b: &Budget{BudgetID: "bud-2"}}
	
	origWithTx := database.WithTx
	defer func() { database.WithTx = origWithTx }()

	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	
	svc := NewServiceWithDB(repo, &sql.DB{})
	
	b, err := svc.CreateBudget(context.Background(), CreateBudgetRequest{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if b.BudgetID != "bud-2" {
		t.Errorf("wrong id")
	}

	repo.err = errors.New("repo fail")
	_, err = svc.CreateBudget(context.Background(), CreateBudgetRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
	repo.err = nil

	database.WithTx = func(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
		return errors.New("tx fail")
	}
	_, err = svc.CreateBudget(context.Background(), CreateBudgetRequest{})
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

	b, err = svc.CreateBudget(context.Background(), CreateBudgetRequest{OwnerID: "123"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if b.BudgetID != "bud-2" {
		t.Errorf("wrong id")
	}
}

type mockAuditRepo struct{}
func (m *mockAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }
func (m *mockAuditRepo) InsertTx(ctx context.Context, tx *sql.Tx, e audit.AuditEvent) error { return nil }

type mockDomainEmitter struct{}
func (m *mockDomainEmitter) EmitTx(ctx context.Context, tx *sql.Tx, e domain.DomainEvent) error { return nil }
