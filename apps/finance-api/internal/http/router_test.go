package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"artha-kosha/apps/finance-api/internal/accounts"
	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/budgets"
	"artha-kosha/apps/finance-api/internal/transactions"
)

type mockAuditRepo struct{}
func (m *mockAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }

type fullMockProvider struct {
	mockAuthProvider
	acc *accounts.Service
	txn *transactions.Service
	bdg *budgets.Service
}

func (m *fullMockProvider) GetAccountsService() *accounts.Service { return m.acc }
func (m *fullMockProvider) GetTransactionsService() *transactions.Service { return m.txn }
func (m *fullMockProvider) GetBudgetsService() *budgets.Service { return m.bdg }
func (m *fullMockProvider) ChangePassword(req auth.ChangePasswordRequest) error { return nil }

func TestNewRouter(t *testing.T) {
	os.Setenv("FINANCE_API_ALLOWED_ORIGINS", "http://localhost:3000")
	p := &fullMockProvider{}
	r := NewRouter(p, &mockAuditRepo{})

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
