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

type mockAuditRepoMalicious struct{}
func (m *mockAuditRepoMalicious) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }

func TestMaliciousCORSRequestHandling(t *testing.T) {
	os.Setenv("FINANCE_API_ALLOWED_ORIGINS", "http://allowed.com")
	
	provider := auth.NewLocalAuthProvider()
	r := router.NewRouter(provider, &mockAuditRepoMalicious{})

	// Malformed request with invalid characters in Origin
	req := httptest.NewRequest("OPTIONS", "/api/v1/test", nil)
	req.Header.Set("Origin", "http://allowed.com\x00")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Origin should not be echoed back since it's malformed/malicious
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected no CORS headers for malformed origin")
	}
	
	os.Unsetenv("FINANCE_API_ALLOWED_ORIGINS")
}
