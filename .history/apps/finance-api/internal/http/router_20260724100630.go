package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	chiCors "github.com/go-chi/cors"

	"artha-kosha/apps/finance-api/internal/auth"
	mw "artha-kosha/apps/finance-api/internal/http/middleware"
)

// NewRouter builds the HTTP router with middleware chain (RequestID → Logging → Recovery → CORS → Timeout → Auth → Router)
func NewRouter(provider auth.AuthProvider) http.Handler {
	r := chi.NewRouter()

	// Middleware chain
	r.Use(chimw.RequestID)
	r.Use(mw.LoggingMiddleware("finance-api"))
	r.Use(chimw.Recoverer)
	r.Use(chiCors.Handler(chiCors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID", "X-Correlation-ID", "X-Session-ID"},
		ExposedHeaders:   []string{"X-Request-ID", "X-Correlation-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(mw.TimeoutMiddleware(30 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Register auth handlers
	RegisterAuthHandlers(r, provider)

	return r
}
 
