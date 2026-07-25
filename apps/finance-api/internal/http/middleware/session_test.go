package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/constants"
	"artha-kosha/apps/finance-api/internal/sessions"
)

type mockAuthSvc struct{}
func (m *mockAuthSvc) Register(req auth.RegisterUserRequest) (auth.RegisterUserResponse, error) { return auth.RegisterUserResponse{}, nil }
func (m *mockAuthSvc) Login(req auth.LoginRequest) (auth.LoginResponse, error) { return auth.LoginResponse{}, nil }
func (m *mockAuthSvc) Logout(sessionID string) error { return nil }
func (m *mockAuthSvc) RevokeAll(userID string) error { return nil }
func (m *mockAuthSvc) ChangePassword(req auth.ChangePasswordRequest) error { return nil }
func (m *mockAuthSvc) DeleteUser(userID string, confirm string, days int) error { return nil }
func (m *mockAuthSvc) GetSession(sessionID string) (sessions.Session, error) {
	if sessionID == "valid" {
		return sessions.Session{ID: "valid", UserID: "u1"}, nil
	}
	if sessionID == "error" {
		return sessions.Session{}, errors.New("db error")
	}
	return sessions.Session{}, errors.New("not found")
}

func TestSessionMiddleware(t *testing.T) {
	h := SessionMiddleware(&mockAuthSvc{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	
	// Valid session
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(constants.HeaderSessionID, "valid")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", w.Code)
	}

	// No session header
	reqInvalid := httptest.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, reqInvalid)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %v", w2.Code)
	}

	// Invalid session header (not found)
	reqNotFound := httptest.NewRequest("GET", "/", nil)
	reqNotFound.Header.Set(constants.HeaderSessionID, "invalid")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, reqNotFound)
	if w3.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %v", w3.Code)
	}
	
	// Error from GetSession
	reqErr := httptest.NewRequest("GET", "/", nil)
	reqErr.Header.Set(constants.HeaderSessionID, "error")
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, reqErr)
	if w4.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %v", w4.Code)
	}
}
