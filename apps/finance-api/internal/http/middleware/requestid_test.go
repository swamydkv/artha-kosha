package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	h := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestID(r.Context())
		corID := GetCorrelationID(r.Context())
		if reqID == "" || corID == "" {
			t.Errorf("expected IDs in context")
		}
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "req-1")
	req.Header.Set("X-Correlation-ID", "cor-1")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
}
