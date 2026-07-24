package auth

import (
	"sync"
	"time"
)

// LoginRateLimiter provides rate limiting for login attempts to mitigate timing attacks and brute forcing
type LoginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewLoginRateLimiter(limit int, window time.Duration) *LoginRateLimiter {
	return &LoginRateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if the given identifier (e.g., username or IP) is allowed to attempt login
func (r *LoginRateLimiter) Allow(identifier string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	times := r.attempts[identifier]

	// filter out old attempts
	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) < r.window {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= r.limit {
		r.attempts[identifier] = validTimes
		return false
	}

	validTimes = append(validTimes, now)
	r.attempts[identifier] = validTimes
	return true
}
