package middleware

import (
    "net/http"

    chiCors "github.com/go-chi/cors"
)

// CorsMiddleware returns a chi CORS handler configured with sensible defaults.
func CorsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    if len(allowedOrigins) == 0 {
        allowedOrigins = []string{"*"}
    }
    o := chiCors.Options{
        AllowedOrigins:   allowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID", "X-Correlation-ID", "X-Session-ID"},
        ExposedHeaders:   []string{"X-Request-ID", "X-Correlation-ID"},
        AllowCredentials: true,
        MaxAge:           300,
    }
    h := chiCors.Handler(o)
    return func(next http.Handler) http.Handler {
        return h(next)
    }
}
