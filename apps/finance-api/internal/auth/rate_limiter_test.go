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

func BenchmarkRateLimiter_Allow(b *testing.B) {
	rl := NewLoginRateLimiter(100, time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow("test_ip")
	}
}
