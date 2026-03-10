package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func HashToken(raw string) value.TokenHash {
	sum := sha256.Sum256([]byte(raw))
	return value.TokenHash(hex.EncodeToString(sum[:]))
}

type TokenHasher struct{}

func NewTokenHasher() *TokenHasher {
	return &TokenHasher{}
}

func (TokenHasher) Hash(raw string) value.TokenHash {
	return HashToken(raw)
}
