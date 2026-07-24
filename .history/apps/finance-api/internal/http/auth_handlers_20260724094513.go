package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"artha-kosha/apps/finance-api/internal/auth"
)

func RegisterAuthHandlers(r chi.Router, provider auth.AuthProvider) {
	r.Post("/register", func(w http.ResponseWriter, req *http.Request) {
		var body auth.RegisterUserRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		result, err := provider.Register(body)
		if err != nil {
			status := http.StatusBadRequest
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				status = http.StatusConflict
			}
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(result)
	})

	r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
		var body auth.LoginRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		result, err := provider.Login(body)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid credentials"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(result)
	})

	r.Post("/logout", func(w http.ResponseWriter, req *http.Request) {
		_ = provider.Logout(strings.TrimSpace(req.Header.Get("X-Session-ID")))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"logged out"}`))
	})

	r.Get("/session", func(w http.ResponseWriter, req *http.Request) {
		sessionID := strings.TrimSpace(req.Header.Get("X-Session-ID"))
		if sessionID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sess, err := provider.GetSession(sessionID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(sess)
	})

	r.Delete("/session/revoke", func(w http.ResponseWriter, req *http.Request) {
		sessionID := strings.TrimSpace(req.Header.Get("X-Session-ID"))
		if sessionID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = provider.Logout(sessionID)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"session revoked"}`))
	})

	r.Delete("/sessions", func(w http.ResponseWriter, req *http.Request) {
		userID := strings.TrimSpace(req.Header.Get("X-User-ID"))
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = provider.RevokeAll(userID)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"all sessions revoked"}`))
	})
}
