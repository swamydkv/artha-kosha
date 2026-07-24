package auth

import (
	"errors"
	"strings"
)

// ChangePasswordRequest represents the request to change a user's password
type ChangePasswordRequest struct {
	Username    string
	OldPassword string
	NewPassword string
}

// ChangePassword updates the user's password if the old password is correct.
// It also ensures that the old password hash is invalidated (by replacing it with a new argon2id hash).
func (p *LocalAuthProvider) ChangePassword(req ChangePasswordRequest) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	user, found := p.users[strings.ToLower(req.Username)]
	if !found {
		return errors.New("user not found")
	}

	match, err := passwordMatches(req.OldPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid old password")
	}

	if len(req.NewPassword) < 12 {
		return errors.New("new password must be at least 12 characters")
	}

	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Invalidate the old hash by replacing it with the newly generated hash
	user.PasswordHash = newHash

	// Note: in a real implementation, we might also want to revoke all existing sessions
	// so the user has to login again with the new password.
	if p.sessSvc != nil {
		_ = p.sessSvc.RevokeAll(user.ID)
	}

	return nil
}
