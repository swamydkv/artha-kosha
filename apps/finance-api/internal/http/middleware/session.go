package middleware

import (
	"context"
	"net/http"

	"artha-kosha/apps/finance-api/internal/auth"
)

type ctxKey int

const (
	SessionCtxKey ctxKey = iota
)

func SessionMiddleware(provider auth.AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("X-Session-ID")
			if sessionID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			sess, err := provider.GetSession(sessionID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), SessionCtxKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
