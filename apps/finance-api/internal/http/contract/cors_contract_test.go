package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mw "artha-kosha/apps/finance-api/internal/http/middleware"
)

func TestCORSPreflightHandling(t *testing.T) {
	// Tests that the CORS middleware handles OPTIONS requests properly
	handler := mw.CorsMiddleware([]string{"*"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for preflight, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers in response")
	}
}
