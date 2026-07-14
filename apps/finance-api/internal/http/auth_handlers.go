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
}
