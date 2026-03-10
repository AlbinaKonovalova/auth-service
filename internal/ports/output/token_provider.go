package output

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type TokenProvider interface {
	// GenerateAccessToken создаёт подписанный JWT из claims.
	GenerateAccessToken(ctx context.Context, claims value.AccessClaims) (token string, expiresIn int64, err error)
	// ParseAccessToken валидирует JWT и возвращает claims.
	ParseAccessToken(ctx context.Context, token string) (value.AccessClaims, error)
	// GenerateRefreshToken генерирует raw refresh token и его hash.
	GenerateRefreshToken(ctx context.Context) (raw string, hash value.TokenHash, err error)
	// GeneratePasswordResetToken генерирует raw password reset token и его hash.
	GeneratePasswordResetToken(ctx context.Context) (raw string, hash value.TokenHash, err error)
}
