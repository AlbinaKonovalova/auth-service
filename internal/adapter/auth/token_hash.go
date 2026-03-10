package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

// HashToken возвращает sha256 hex digest от raw token строки.
func HashToken(raw string) value.TokenHash {
	sum := sha256.Sum256([]byte(raw))
	return value.TokenHash(hex.EncodeToString(sum[:]))
}

// TokenHasher реализует output.TokenHasher.
type TokenHasher struct{}

func NewTokenHasher() *TokenHasher {
	return &TokenHasher{}
}

func (TokenHasher) Hash(raw string) value.TokenHash {
	return HashToken(raw)
}
