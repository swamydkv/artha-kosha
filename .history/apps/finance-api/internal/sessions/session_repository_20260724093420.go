package sessions
package sessions

import (
    "errors"
    "sync"
    "time"
)

type Status string

const (
    StatusActive  Status = "active"
    StatusExpired Status = "expired"
    StatusRevoked Status = "revoked"
)

type Session struct {
    ID             string    `json:"id"`
    UserID         string    `json:"user_id"`
    CreatedAt      time.Time `json:"created_at"`
    LastActivityAt time.Time `json:"last_activity_at"`
    ExpiresAt      time.Time `json:"expires_at"`
    RevokedAt      time.Time `json:"revoked_at"`
    UserAgent      string    `json:"user_agent"`
    IPAddress      string    `json:"ip_address"`
    Status         Status    `json:"status"`
}

type Repository interface {
    Create(s Session) error
    Get(id string) (Session, error)
    Revoke(id string) error
    RevokeAllByUser(userID string) error
    UpdateActivity(id string, lastActivity time.Time) error
}

type InMemoryRepo struct {
    mu       sync.RWMutex
    sessions map[string]Session
}

func NewInMemoryRepo() *InMemoryRepo {
    return &InMemoryRepo{sessions: make(map[string]Session)}
}

func (r *InMemoryRepo) Create(s Session) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if s.ID == "" {
        return errors.New("session id required")
    }
    r.sessions[s.ID] = s
    return nil
}

func (r *InMemoryRepo) Get(id string) (Session, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, ok := r.sessions[id]
    if !ok {
        return Session{}, errors.New("not found")
    }
    // expire if past ExpiresAt
    if !s.ExpiresAt.IsZero() && time.Now().After(s.ExpiresAt) {
        s.Status = StatusExpired
        r.sessions[id] = s
    }
    return s, nil
}

func (r *InMemoryRepo) Revoke(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    s, ok := r.sessions[id]
    if !ok {
        return nil
    }
    now := time.Now()
    s.Status = StatusRevoked
    s.RevokedAt = now
    r.sessions[id] = s
    return nil
}

func (r *InMemoryRepo) RevokeAllByUser(userID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now()
    for id, s := range r.sessions {
        if s.UserID == userID {
            s.Status = StatusRevoked
            s.RevokedAt = now
            r.sessions[id] = s
        }
    }
    return nil
}

func (r *InMemoryRepo) UpdateActivity(id string, lastActivity time.Time) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    s, ok := r.sessions[id]
    if !ok {
        return errors.New("not found")
    }
    s.LastActivityAt = lastActivity
    r.sessions[id] = s
    return nil
}
