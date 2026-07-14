package http

import (
	"encoding/json"
	"net/http"

	"artha-kosha/apps/finance-api/internal/auth"
)

// RegisterHandler handles registration requests
type RegisterHandler struct {
	registerService *auth.RegisterService
}

// NewRegisterHandler creates a new registration handler
func NewRegisterHandler(registerService *auth.RegisterService) *RegisterHandler {
	return &RegisterHandler{
		registerService: registerService,
	}
}

// ServeHTTP handles the HTTP request for registration
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.registerService.Register(r.Context(), req)
	if err != nil {
		// Determine status code based on error
		statusCode := http.StatusBadRequest
		if err.Error() == "username already exists" || 
		   err.Error() == "email already exists" || 
		   err.Error() == "mobile number already exists" {
			statusCode = http.StatusConflict
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}