package middleware

import (
	"context"
	"net/http"
	"github.com/google/uuid"

	"artha-kosha/apps/finance-api/internal/constants"
)

type requestIDKey struct{}
type correlationIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func GetCorrelationID(ctx context.Context) string {
	if v := ctx.Value(correlationIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RequestIDMiddleware ensures every request has X-Request-ID and X-Correlation-ID headers
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(constants.HeaderRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		corrID := r.Header.Get(constants.HeaderCorrelationID)
		if corrID == "" {
			corrID = reqID
		}

		// set headers on response for clients
		w.Header().Set(constants.HeaderRequestID, reqID)
		w.Header().Set(constants.HeaderCorrelationID, corrID)

		ctx := r.Context()
		ctx = WithRequestID(ctx, reqID)
		ctx = WithCorrelationID(ctx, corrID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
