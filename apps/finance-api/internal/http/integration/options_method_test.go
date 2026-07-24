package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/auth"
	router "artha-kosha/apps/finance-api/internal/http"
)

type mockAuditRepoOptions struct{}

func (m *mockAuditRepoOptions) Insert(ctx context.Context, e audit.AuditEvent) error { return nil }

func TestOptionsMethodHandling(t *testing.T) {
	provider := auth.NewLocalAuthProvider()
	r := router.NewRouter(provider, &mockAuditRepoOptions{})

	req := httptest.NewRequest("OPTIONS", "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for OPTIONS preflight, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods in response")
	}
}
