package middleware

import (
    "context"
    "net/http"
    "time"
)

func TimeoutMiddleware(d time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), d)
            defer cancel()
            r = r.WithContext(ctx)
            done := make(chan struct{})
            go func() {
                next.ServeHTTP(w, r)
                close(done)
            }()
            select {
            case <-ctx.Done():
                if ctx.Err() == context.DeadlineExceeded {
                    http.Error(w, "request timeout", http.StatusGatewayTimeout)
                }
            case <-done:
                return
            }
        })
    }
}
