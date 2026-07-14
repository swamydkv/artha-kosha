package http

import (
	"net/http"

	"artha-kosha/apps/finance-api/internal/auth"
)

// addCORSHeaders adds CORS headers to the response
func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Correlation-ID, X-Session-ID")
	w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Correlation-ID")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "300")
}

func NewRouter(provider auth.AuthProvider) *http.ServeMux {
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w)
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	
	// Register auth handlers
	RegisterAuthHandlers(mux, provider)
	
	return mux
}
