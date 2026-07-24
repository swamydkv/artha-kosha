package sessions

import (
    "time"
)

type Service struct {
    repo Repository
    // default TTL for sessions
    ttl time.Duration
}

func NewService(repo Repository, ttl time.Duration) *Service {
    return &Service{repo: repo, ttl: ttl}
}

func (s *Service) CreateSession(id string, userID string, userAgent string, ip string) (Session, error) {
    now := time.Now().UTC()
    sess := Session{
        ID:             id,
        UserID:         userID,
        CreatedAt:      now,
        LastActivityAt: now,
        ExpiresAt:      now.Add(s.ttl),
        UserAgent:      userAgent,
        IPAddress:      ip,
        Status:         StatusActive,
    }
    if err := s.repo.Create(sess); err != nil {
        return Session{}, err
    }
    return sess, nil
}

func (s *Service) GetSession(id string) (Session, error) {
    return s.repo.Get(id)
}

func (s *Service) RevokeSession(id string) error {
    return s.repo.Revoke(id)
}

func (s *Service) RevokeAll(userID string) error {
    return s.repo.RevokeAllByUser(userID)
}

func (s *Service) RefreshActivity(id string) error {
    return s.repo.UpdateActivity(id, time.Now().UTC())
}
