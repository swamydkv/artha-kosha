package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/audit"
	router "artha-kosha/apps/finance-api/internal/http"
)

type mockAuditRepoSecurity struct{}
func (m *mockAuditRepoSecurity) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }

func TestUnauthorizedOriginRejection(t *testing.T) {
	// Set specific allowed origin
	os.Setenv("FINANCE_API_ALLOWED_ORIGINS", "http://allowed.com")
	
	provider := auth.NewLocalAuthProvider()
	r := router.NewRouter(provider, &mockAuditRepoSecurity{})

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("Origin", "http://evil-hacker.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Origin should not be echoed back
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected no CORS headers for unauthorized origin")
	}
	
	// Reset env
	os.Unsetenv("FINANCE_API_ALLOWED_ORIGINS")
}
