package sessions

import (
	"testing"
	"time"
)

func TestService_CreateSession(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	sess, err := svc.CreateSession("1", "u1", "agent1", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sess.UserID != "u1" || sess.IPAddress != "127.0.0.1" || sess.UserAgent != "agent1" {
		t.Errorf("session values mismatch")
	}
}

func TestService_GetSession(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	sess, _ := svc.CreateSession("1", "u1", "agent1", "127.0.0.1")
	
	got, err := svc.GetSession(sess.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != sess.ID {
		t.Errorf("expected %s, got %s", sess.ID, got.ID)
	}

	// Revoke and test
	svc.RevokeSession(sess.ID)
	_, err = svc.GetSession(sess.ID)
	if err == nil {
		t.Error("expected error for revoked session")
	}

	// Expired session
	sess2, _ := svc.CreateSession("2", "u2", "", "")
	sess2FromRepo, _ := repo.Get(sess2.ID)
	sess2FromRepo.ExpiresAt = time.Now().Add(-time.Hour)
	repo.Create(sess2FromRepo)

	_, err = svc.GetSession(sess2.ID)
	if err == nil || err.Error() != "session is not active" {
		t.Errorf("expected session is not active error, got %v", err)
	}
}

func TestService_RevokeSession(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	sess, _ := svc.CreateSession("1", "u1", "", "")
	err := svc.RevokeSession(sess.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.Get(sess.ID)
	if got.Status != StatusRevoked {
		t.Errorf("expected revoked, got %s", got.Status)
	}
}

func TestService_RevokeAllSessions(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	svc.CreateSession("1", "u1", "", "")
	svc.CreateSession("2", "u1", "", "")
	svc.CreateSession("3", "u2", "", "")

	err := svc.RevokeAll("u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s1, _ := repo.Get("1")
	if s1.Status != StatusRevoked {
		t.Errorf("expected revoked for 1")
	}

	s3, _ := repo.Get("3")
	if s3.Status == StatusRevoked {
		t.Errorf("did not expect revoked for 3")
	}
}

func TestService_RefreshActivity(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	sess, _ := svc.CreateSession("1", "u1", "", "")
	err := svc.RefreshActivity(sess.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_Errors(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Hour)

	// Create error
	_, err := svc.CreateSession("", "u1", "", "") // empty ID causes Create error in InMemoryRepo
	if err == nil {
		t.Error("expected error")
	}

	// Get error
	_, err = svc.GetSession("nonexistent")
	if err == nil {
		t.Error("expected error")
	}

	// Expired check is covered in previous test but we can also mock repo returning active but past expires_at
	// Wait, InMemoryRepo mutates status, so it's already covered.
}
