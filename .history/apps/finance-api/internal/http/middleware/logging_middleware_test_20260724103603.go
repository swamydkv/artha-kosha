package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestLoggingMiddleware_InsertsLoggerIntoContext(t *testing.T) {
    called := false
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        l := GetLogger(r.Context())
        if l == nil {
            t.Fatal("expected logger in context")
        }
        called = true
        w.WriteHeader(http.StatusOK)
    })

    handler := LoggingMiddleware("test-service")(next)

    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if !called {
        t.Fatal("next handler was not called")
    }
}
