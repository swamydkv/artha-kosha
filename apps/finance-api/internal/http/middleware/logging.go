package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"artha-kosha/apps/finance-api/internal/constants"
	"log/slog"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}
func GetLogger(ctx context.Context) *slog.Logger {
	if v := ctx.Value(loggerKey{}); v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return slog.Default()
}

func LoggingMiddleware(service string) func(http.Handler) http.Handler {
	// create base logger
	l := slog.New(slog.NewJSONHandler(log.Default().Writer(), nil))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := r.Header.Get(constants.HeaderRequestID)
			corrID := r.Header.Get(constants.HeaderCorrelationID)
			ctx := WithLogger(r.Context(), l)
			// try to extract user and session from context
			userID := ""
			sessionID := ""
			if sess, ok := r.Context().Value("session_id").(string); ok {
				sessionID = sess
			}
			if uid, ok := r.Context().Value("user_id").(string); ok {
				userID = uid
			}

			// add basic fields
			l.InfoContext(ctx, "request.started",
				"service", service,
				"component", "http",
				"operation", "request.started",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", reqID,
				"correlation_id", corrID,
				"user_id", userID,
				"session_id", sessionID,
			)
			next.ServeHTTP(w, r.WithContext(ctx))

			// get possibly updated values
			if sess, ok := r.Context().Value("session_id").(string); ok {
				sessionID = sess
			}
			if uid, ok := r.Context().Value("user_id").(string); ok {
				userID = uid
			}

			l.InfoContext(ctx, "request.finished",
				"service", service,
				"component", "http",
				"operation", "request.finished",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", reqID,
				"correlation_id", corrID,
				"user_id", userID,
				"session_id", sessionID,
				"duration", time.Since(start).Milliseconds(),
			)
		})
	}
}
