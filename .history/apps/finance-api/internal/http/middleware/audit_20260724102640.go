package middleware

import (
    "context"
    "net/http"

    "artha-kosha/apps/finance-api/internal/audit"
)

func AuditMiddleware(auditRepo audit.Repository) func(http.Handler) http.Handler {
    if auditRepo == nil {
        return func(h http.Handler) http.Handler { return h }
    }
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // basic audit: record request path and method
            ctx := r.Context()
            // fire and forget; don't block request on audit failure
            go func() {
                _ = auditRepo.Insert(ctx, audit.AuditEvent{
                    ID:        "",
                    RequestID: r.Header.Get("X-Request-ID"),
                    UserID:    "",
                    SessionID: r.Header.Get("X-Session-ID"),
                    Resource:  r.URL.Path,
                    Action:    r.Method,
                    Result:    "success",
                    Timestamp:  r.Context().Value("time").(time.Time),
                })
            }()
            next.ServeHTTP(w, r)
        })
    }
}
