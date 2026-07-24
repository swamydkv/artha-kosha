package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/auth"
	mw "artha-kosha/apps/finance-api/internal/http/middleware"
)

// NewRouter builds the HTTP router with middleware chain (RequestID → Logging → Recovery → CORS → Timeout → Audit → Auth → Router)
func NewRouter(provider auth.AuthProvider, auditRepo audit.Repository) http.Handler {
	r := chi.NewRouter()

	// Middleware chain
	r.Use(chimw.RequestID)
	r.Use(mw.LoggingMiddleware("finance-api"))
	r.Use(chimw.Recoverer)
	r.Use(mw.CorsMiddleware(nil))
	r.Use(mw.TimeoutMiddleware(30 * time.Second))
	r.Use(mw.AuditMiddleware(auditRepo))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Register auth handlers
	RegisterAuthHandlers(r, provider)

	return r
}
