package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}

func passwordMatches(password string, hash string) bool {
	return hashPassword(password) == hash
}
