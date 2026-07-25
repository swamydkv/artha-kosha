package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/users"
)

type mockUserRepo struct {
	err error
}
func (m *mockUserRepo) CreateUser(ctx context.Context, u users.CreateUserRequest) (*users.User, error) { return nil, m.err }
func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (*users.User, error) { return nil, m.err }
func (m *mockUserRepo) GetUserByUsername(ctx context.Context, un string) (*users.User, error) { return nil, m.err }
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, em string) (*users.User, error) { return nil, m.err }
func (m *mockUserRepo) GetUserByMobileNumber(ctx context.Context, mb string) (*users.User, error) { return nil, m.err }
func (m *mockUserRepo) CheckUserExists(ctx context.Context, username, email, mobile string) (*users.UserExistsCheck, error) { return nil, m.err }
func (m *mockUserRepo) CreateSession(ctx context.Context, sID, uID, ip, ua string) error { return m.err }
func (m *mockUserRepo) DeleteUser(ctx context.Context, id string, r int) error { return m.err }

func TestRegisterHandler_ServeHTTP(t *testing.T) {
	svc := auth.NewRegisterService(nil, &mockUserRepo{})
	h := NewRegisterHandler(svc)
	
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	req2 := httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{}`))
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	// validation error
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w2.Code)
	}
}
