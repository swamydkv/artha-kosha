package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware_AllowsOriginAndPreflight(t *testing.T) {
	handler := CorsMiddleware([]string{"https://example.com"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %d", rr.Code)
	}

	got := rr.Header().Get("Access-Control-Allow-Origin")
	if got != "https://example.com" {
		t.Fatalf("expected Access-Control-Allow-Origin 'https://example.com' but got '%s'", got)
	}
}

func TestCorsMiddleware_EmptyOrigins(t *testing.T) {
	handler := CorsMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "https://other.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %d", rr.Code)
	}

	got := rr.Header().Get("Access-Control-Allow-Origin")
	if got != "*" {
		t.Fatalf("expected Access-Control-Allow-Origin '*' but got '%s'", got)
	}
}
