package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken returns SHA256 hash of token (for storage).
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
