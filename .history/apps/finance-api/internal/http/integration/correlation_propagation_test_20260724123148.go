package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "artha-kosha/apps/finance-api/internal/http/middleware"
)

func TestCorrelationPropagation(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ensure headers are present on response
        reqID := w.Header().Get("X-Request-ID")
        corrID := w.Header().Get("X-Correlation-ID")
        if reqID == "" || corrID == "" {
            t.Fatalf("missing propagated headers: req=%q corr=%q", reqID, corrID)
        }
        w.WriteHeader(200)
    })

    h := middleware.RequestIDMiddleware(middleware.LoggingMiddleware("svc")(handler))
    req := httptest.NewRequest("POST", "/test", nil)
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)
    if rr.Code != 200 {
        t.Fatalf("unexpected status: %d", rr.Code)
    }
    if rr.Header().Get("X-Request-ID") == "" || rr.Header().Get("X-Correlation-ID") == "" {
        t.Fatalf("response headers not set")
    }
}
