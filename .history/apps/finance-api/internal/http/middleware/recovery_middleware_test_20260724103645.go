package middleware

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestRecoveryMiddleware_RecoversAndReturns500(t *testing.T) {
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        panic("boom")
    })

    handler := RecoveryMiddleware()(next)

    req := httptest.NewRequest("GET", "/", nil)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected status %d but got %d", http.StatusInternalServerError, rr.Code)
    }

    if !strings.Contains(rr.Body.String(), "internal server error") {
        t.Fatalf("expected internal server error body but got %q", rr.Body.String())
    }
}
