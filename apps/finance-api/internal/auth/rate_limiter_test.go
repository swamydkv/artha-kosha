package auth

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := NewLoginRateLimiter(2, time.Minute)
	
	if !rl.Allow("ip1") {
		t.Error("expected allow")
	}
	if !rl.Allow("ip1") {
		t.Error("expected allow")
	}
	if rl.Allow("ip1") {
		t.Error("expected reject")
	}
}
