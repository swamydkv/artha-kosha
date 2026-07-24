package validation

import (
	"testing"
	"time"

	"artha-kosha/apps/finance-api/internal/sessions"
)

func TestSessionValidation_T032(t *testing.T) {
	repo := sessions.NewInMemoryRepo()
	svc := sessions.NewService(repo, 1*time.Hour)

	// Test Create
	sess, err := svc.CreateSession("sess-1", "user-1", "Go-http-client/1.1", "127.0.0.1")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	if sess.ID != "sess-1" {
		t.Errorf("expected sess-1, got %s", sess.ID)
	}

	// Test Get
	fetched, err := svc.GetSession("sess-1")
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}
	if fetched.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", fetched.UserID)
	}

	// Test Revoke
	err = svc.RevokeSession("sess-1")
	if err != nil {
		t.Fatalf("failed to revoke: %v", err)
	}

	// Get after revoke should still work for in memory repo but status might be revoked.
	// We just test repo level validation here.
	_, err = svc.GetSession("sess-1")
	if err == nil {
		t.Errorf("expected error when getting revoked session, got nil")
	}
}
