package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"artha-kosha/apps/finance-api/internal/audit"
	"artha-kosha/apps/finance-api/internal/constants"
)

func genID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return prefix + "-" + hex.EncodeToString(buf)
}

func AuditMiddleware(auditRepo audit.Repository) func(http.Handler) http.Handler {
	if auditRepo == nil {
		return func(h http.Handler) http.Handler { return h }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ua := r.UserAgent()
			clientIP := r.RemoteAddr
			// try common forwarded header
			if xf := r.Header.Get(constants.HeaderForwardedFor); xf != "" {
				clientIP = strings.Split(xf, ",")[0]
			}

			// No payload capture currently
			// fire and forget; don't block request on audit failure
			go func() {
				_ = auditRepo.Insert(ctx, audit.AuditEvent{
					ID:        genID("audit"),
					RequestID: r.Header.Get(constants.HeaderRequestID),
					UserID:    "", // usually extracted from session
					SessionID: r.Header.Get(constants.HeaderSessionID),
					Resource:  r.URL.Path,
					Action:    r.Method,
					Result:    "success",
					Timestamp: time.Now().UTC(),
					UserAgent: ua,
					ClientIP:  clientIP,
				})
			}()
			next.ServeHTTP(w, r)
		})
	}
}
