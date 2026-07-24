package http

import (
	"encoding/json"
	nethttp "net/http"
	"strings"

	"artha-kosha/apps/finance-api/internal/auth"
)

func RegisterAuthHandlers(mux *nethttp.ServeMux, provider auth.AuthProvider) {
	mux.HandleFunc("/register", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)
		
		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}
		
		if r.Method != nethttp.MethodPost {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		var req auth.RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		result, err := provider.Register(req)
		if err != nil {
			status := nethttp.StatusBadRequest
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				status = nethttp.StatusConflict
			}
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusCreated)
		_ = json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/login", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)
		
		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}
		
		if r.Method != nethttp.MethodPost {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		var req auth.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		result, err := provider.Login(req)
		if err != nil {
			w.WriteHeader(nethttp.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid credentials"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/logout", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)
		
		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}
		
		if r.Method != nethttp.MethodPost {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}
		_ = provider.Logout(strings.TrimSpace(r.Header.Get("X-Session-ID")))
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte(`{"message":"logged out"}`))
	})

	// GET /session - validate current session
	mux.HandleFunc("/session", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)

		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}

		if r.Method != nethttp.MethodGet {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		sessionID := strings.TrimSpace(r.Header.Get("X-Session-ID"))
		if sessionID == "" {
			w.WriteHeader(nethttp.StatusUnauthorized)
			return
		}

		sess, err := provider.GetSession(sessionID)
		if err != nil {
			w.WriteHeader(nethttp.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(sess)
	})

	// DELETE /session - revoke current session
	mux.HandleFunc("/session/revoke", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)

		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}

		if r.Method != nethttp.MethodDelete {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		sessionID := strings.TrimSpace(r.Header.Get("X-Session-ID"))
		if sessionID == "" {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		_ = provider.Logout(sessionID)
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte(`{"message":"session revoked"}`))
	})

	// DELETE /sessions - revoke all sessions for a user
	mux.HandleFunc("/sessions", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		addCORSHeaders(w)

		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}

		if r.Method != nethttp.MethodDelete {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
		if userID == "" {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		_ = provider.RevokeAll(userID)
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte(`{"message":"all sessions revoked"}`))
	})
}
