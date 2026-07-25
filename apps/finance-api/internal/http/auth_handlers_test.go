package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/sessions"
)

type mockAuthProvider struct {
	err error
}
func (m *mockAuthProvider) Register(req auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	if req.Username == "duplicate" { return auth.RegisterUserResponse{}, errors.New("duplicate key") }
	if m.err != nil { return auth.RegisterUserResponse{}, m.err }
	return auth.RegisterUserResponse{}, nil
}
func (m *mockAuthProvider) Login(req auth.LoginRequest) (auth.LoginResponse, error) {
	if m.err != nil { return auth.LoginResponse{}, m.err }
	return auth.LoginResponse{}, nil
}
func (m *mockAuthProvider) GetSession(sessionID string) (sessions.Session, error) {
	if m.err != nil { return sessions.Session{}, m.err }
	return sessions.Session{}, nil
}
func (m *mockAuthProvider) Logout(sessionID string) error { return m.err }
func (m *mockAuthProvider) RevokeAll(userID string) error { return m.err }
func (m *mockAuthProvider) DeleteUser(userID, confirmation string, retentionDays int) error {
	if confirmation != "DELETE" { return errors.New("confirmation must be exactly 'DELETE'") }
	return m.err
}
func (m *mockAuthProvider) ChangePassword(req auth.ChangePasswordRequest) error { return m.err }

func TestAuthHandlers(t *testing.T) {
	r := chi.NewRouter()
	p := &mockAuthProvider{}
	RegisterAuthHandlers(r, p)

	// 1. Register - invalid body
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("{invalid"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 2. Register - duplicate
	req = httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"username":"duplicate"}`))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict { t.Errorf("expected 409") }

	// 3. Register - provider error
	p.err = errors.New("err")
	req = httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"username":"test"}`))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }
	p.err = nil

	// 4. Register - success
	req = httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"username":"test"}`))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated { t.Errorf("expected 201") }

	// 5. Login - invalid body
	req = httptest.NewRequest("POST", "/login", bytes.NewBufferString("{invalid"))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 6. Login - provider error
	p.err = errors.New("err")
	req = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"test"}`))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized { t.Errorf("expected 401") }
	p.err = nil

	// 7. Login - success
	req = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"test"}`))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }

	// 8. Logout
	req = httptest.NewRequest("POST", "/logout", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }

	// 9. Session - no header
	req = httptest.NewRequest("GET", "/session", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized { t.Errorf("expected 401") }

	// 10. Session - provider error
	p.err = errors.New("err")
	req = httptest.NewRequest("GET", "/session", nil)
	req.Header.Set("X-Session-ID", "sess-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized { t.Errorf("expected 401") }
	p.err = nil

	// 11. Session - success
	req = httptest.NewRequest("GET", "/session", nil)
	req.Header.Set("X-Session-ID", "sess-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }

	// 12. Session revoke - no header
	req = httptest.NewRequest("DELETE", "/session/revoke", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 13. Session revoke - success
	req = httptest.NewRequest("DELETE", "/session/revoke", nil)
	req.Header.Set("X-Session-ID", "sess-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }

	// 14. Sessions revoke all - no header
	req = httptest.NewRequest("DELETE", "/sessions", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 15. Sessions revoke all - success
	req = httptest.NewRequest("DELETE", "/sessions", nil)
	req.Header.Set("X-User-ID", "usr-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }

	// 16. Delete user - no header
	req = httptest.NewRequest("DELETE", "/user/account", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized { t.Errorf("expected 401") }

	// 17. Delete user - invalid body
	req = httptest.NewRequest("DELETE", "/user/account", bytes.NewBufferString("{invalid"))
	req.Header.Set("X-User-ID", "usr-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 18. Delete user - provider error (wrong confirmation)
	req = httptest.NewRequest("DELETE", "/user/account", bytes.NewBufferString(`{"confirmation":"WRONG"}`))
	req.Header.Set("X-User-ID", "usr-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest { t.Errorf("expected 400") }

	// 19. Delete user - provider error (other error)
	p.err = errors.New("err")
	req = httptest.NewRequest("DELETE", "/user/account", bytes.NewBufferString(`{"confirmation":"DELETE"}`))
	req.Header.Set("X-User-ID", "usr-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError { t.Errorf("expected 500") }
	p.err = nil

	// 20. Delete user - success
	req = httptest.NewRequest("DELETE", "/user/account", bytes.NewBufferString(`{"confirmation":"DELETE"}`))
	req.Header.Set("X-User-ID", "usr-1")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("expected 200") }
}
