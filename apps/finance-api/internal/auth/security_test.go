package auth

import (
	"testing"
	"time"
)

func TestHashAndCompare_Success(t *testing.T) {
	pw := "Str0ngP@ssw0rd!"
	h, err := hashPassword(pw)
	if err != nil {
		t.Fatalf("hashPassword error: %v", err)
	}
	ok, err := passwordMatches(pw, h)
	if err != nil {
		t.Fatalf("passwordMatches error: %v", err)
	}
	if !ok {
		t.Fatalf("expected password to match generated hash")
	}
}

func TestPasswordMatches_FailureAndTiming(t *testing.T) {
	pw := "Str0ngP@ssw0rd!"
	h, err := hashPassword(pw)
	if err != nil {
		t.Fatalf("hashPassword error: %v", err)
	}
	// wrong password should not match
	ok, err := passwordMatches("wrong-password", h)
	if err != nil {
		t.Fatalf("passwordMatches error: %v", err)
	}
	if ok {
		t.Fatalf("expected wrong password to not match")
	}

	// basic timing sanity: both calls return within reasonable time
	start := time.Now()
	_, _ = passwordMatches(pw, h)
	_, _ = passwordMatches("wrong-password", h)
	if time.Since(start) > 5*time.Second {
		t.Fatalf("passwordMatches calls took too long")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	pw := "Str0ngP@ssw0rd!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hashPassword(pw)
	}
}

func BenchmarkPasswordMatches(b *testing.B) {
	pw := "Str0ngP@ssw0rd!"
	h, _ := hashPassword(pw)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = passwordMatches(pw, h)
	}
}
