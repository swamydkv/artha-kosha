package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/audit"
	router "artha-kosha/apps/finance-api/internal/http"
)

type mockAuditRepo struct{}
func (m *mockAuditRepo) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }

func TestCORSHeadersIntegration(t *testing.T) {
	provider := auth.NewLocalAuthProvider()
	r := router.NewRouter(provider, &mockAuditRepo{})

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("Origin", "http://allowed-origin.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Note: go-chi/cors returns the origin directly if it's allowed
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers in response")
	}
}
