package http
package http

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "artha-kosha/apps/finance-api/internal/auth"
)

func TestNewRouter_HealthEndpoint(t *testing.T) {
    provider := auth.NewLocalAuthProvider()
    r := NewRouter(provider, nil)

    req := httptest.NewRequest("GET", "/health", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status 200 but got %d", rr.Code)
    }
}
