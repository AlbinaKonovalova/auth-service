package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type JWTClaims struct {
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type TokenProviderConfig struct {
	Secret            string
	AccessTokenTTL    time.Duration
	RefreshTokenBytes int
}

type TokenProvider struct {
	cfg TokenProviderConfig
}

func NewTokenProvider(cfg TokenProviderConfig) *TokenProvider {
	return &TokenProvider{cfg: cfg}
}

func (p *TokenProvider) GenerateAccessToken(_ context.Context, claims value.AccessClaims) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(p.cfg.AccessTokenTTL)

	jwtClaims := JWTClaims{
		Email:       claims.Email,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signed, err := token.SignedString([]byte(p.cfg.Secret))
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return signed, int64(p.cfg.AccessTokenTTL.Seconds()), nil
}

func (p *TokenProvider) ParseAccessToken(_ context.Context, tokenStr string) (value.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(p.cfg.Secret), nil
	})
	if err != nil {
		return value.AccessClaims{}, fmt.Errorf("parse token: %w", err)
	}

	c, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return value.AccessClaims{}, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(c.Subject)
	if err != nil {
		return value.AccessClaims{}, fmt.Errorf("parse user id: %w", err)
	}

	return value.AccessClaims{
		UserID:      userID,
		Email:       c.Email,
		Roles:       c.Roles,
		Permissions: c.Permissions,
	}, nil
}

func (p *TokenProvider) GenerateRefreshToken(_ context.Context) (string, value.TokenHash, error) {
	return p.generateSecureToken()
}

func (p *TokenProvider) GeneratePasswordResetToken(_ context.Context) (string, value.TokenHash, error) {
	return p.generateSecureToken()
}

func (p *TokenProvider) generateSecureToken() (string, value.TokenHash, error) {
	size := p.cfg.RefreshTokenBytes
	if size == 0 {
		size = 32
	}

	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate token: %w", err)
	}

	raw := hex.EncodeToString(b)
	hash := HashToken(raw)

	return raw, hash, nil
}
