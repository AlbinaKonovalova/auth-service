package output

import "github.com/AlbinaKonovalova/auth-service/internal/domain/value"

type TokenHasher interface {
	Hash(raw string) value.TokenHash
}
