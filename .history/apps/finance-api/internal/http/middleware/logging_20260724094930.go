package middleware

import (
    "context"
    "log"
    "net/http"
    "time"

    "log/slog"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context { return context.WithValue(ctx, loggerKey{}, l) }
func GetLogger(ctx context.Context) *slog.Logger {
    if v := ctx.Value(loggerKey{}); v != nil {
        if l, ok := v.(*slog.Logger); ok { return l }
    }
    return slog.Default()
}

func LoggingMiddleware(service string) func(http.Handler) http.Handler {
    // create base logger
    l := slog.New(slog.NewJSONHandler(log.Default().Writer(), nil))
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            reqID := r.Header.Get("X-Request-ID")
            ctx := WithLogger(r.Context(), l)
            // add basic fields
            l.InfoContext(ctx, "request.started", "method", r.Method, "path", r.URL.Path, "request_id", reqID)
            next.ServeHTTP(w, r.WithContext(ctx))
            l.InfoContext(ctx, "request.finished", "method", r.Method, "path", r.URL.Path, "request_id", reqID, "duration_ms", time.Since(start).Milliseconds())
        })
    }
}
