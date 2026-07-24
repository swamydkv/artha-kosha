package auth

import (
	"crypto/subtle"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func hashPassword(password string) (string, error) {
	// Use reasonable defaults; can be tuned later
	params := &argon2id.Params{Memory: 64 * 1024, Iterations: 3, Parallelism: 2, SaltLength: 16, KeyLength: 32}
	s, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", fmt.Errorf("hash error: %w", err)
	}
	return s, nil
}

func passwordMatches(password string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	// Constant time check of bool
	var v byte = 0
	if match {
		v = 1
	}
	return subtle.ConstantTimeCompare([]byte{v}, []byte{1}) == 1, nil
}
