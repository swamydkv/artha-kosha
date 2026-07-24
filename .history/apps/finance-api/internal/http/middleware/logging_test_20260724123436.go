package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware_SetsLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := GetLogger(r.Context())
		if l == nil {
			t.Fatalf("logger not set in context")
		}
		w.WriteHeader(200)
	})

	h := RequestIDMiddleware(LoggingMiddleware("test-service")(handler))
	req := httptest.NewRequest("GET", "/foo", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("unexpected status: %d", rr.Code)
	}
}
