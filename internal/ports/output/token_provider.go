package output

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type TokenProvider interface {
	GenerateAccessToken(ctx context.Context, claims value.AccessClaims) (token string, expiresIn int64, err error)
	ParseAccessToken(ctx context.Context, token string) (value.AccessClaims, error)
	GenerateRefreshToken(ctx context.Context) (raw string, hash value.TokenHash, err error)
	GeneratePasswordResetToken(ctx context.Context) (raw string, hash value.TokenHash, err error)
}
