package middleware

import (
	"net/http"

	"artha-kosha/apps/finance-api/internal/constants"

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
		AllowedHeaders:   []string{"Content-Type", "Authorization", constants.HeaderRequestID, constants.HeaderCorrelationID, constants.HeaderSessionID, constants.HeaderUserID},
		ExposedHeaders:   []string{constants.HeaderRequestID, constants.HeaderCorrelationID},
		AllowCredentials: true,
		MaxAge:           300,
	}
	h := chiCors.Handler(o)
	return func(next http.Handler) http.Handler {
		return h(next)
	}
}
