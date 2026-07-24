package sessions

import (
	"testing"
	"time"
)

func TestInMemoryRepo_CreateAndGet(t *testing.T) {
	repo := NewInMemoryRepo()
	s := Session{
		ID:        "1",
		UserID:    "u1",
		Status:    StatusActive,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := repo.Create(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != "u1" {
		t.Errorf("expected u1, got %s", got.UserID)
	}

	_, err = repo.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}

	// test empty ID
	err = repo.Create(Session{})
	if err == nil {
		t.Error("expected error for empty session id")
	}
}

func TestInMemoryRepo_Revoke(t *testing.T) {
	repo := NewInMemoryRepo()
	s := Session{ID: "1", Status: StatusActive}
	repo.Create(s)

	err := repo.Revoke("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.Get("1")
	if got.Status != StatusRevoked {
		t.Errorf("expected revoked, got %s", got.Status)
	}

	// Revoke nonexistent
	repo.Revoke("nonexistent")
}

func TestInMemoryRepo_RevokeAllByUser(t *testing.T) {
	repo := NewInMemoryRepo()
	repo.Create(Session{ID: "1", UserID: "u1"})
	repo.Create(Session{ID: "2", UserID: "u1"})
	repo.Create(Session{ID: "3", UserID: "u2"})

	err := repo.RevokeAllByUser("u1")
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

func TestInMemoryRepo_UpdateActivity(t *testing.T) {
	repo := NewInMemoryRepo()
	s := Session{ID: "1", Status: StatusActive}
	repo.Create(s)

	now := time.Now()
	err := repo.UpdateActivity("1", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.Get("1")
	if !got.LastActivityAt.Equal(now) {
		t.Errorf("expected activity updated")
	}

	err = repo.UpdateActivity("nonexistent", now)
	if err == nil {
		t.Error("expected error for nonexistent")
	}
}
