package entity

import (
	"strings"
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash value.TokenHash
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func NewRefreshSession(
	id uuid.UUID,
	userID uuid.UUID,
	tokenHash value.TokenHash,
	now time.Time,
	ttl time.Duration,
) (RefreshSession, error) {
	if id == uuid.Nil {
		return RefreshSession{}, domain.ErrInvalidRefreshSessionID
	}
	if userID == uuid.Nil {
		return RefreshSession{}, domain.ErrInvalidUserID
	}
	if strings.TrimSpace(tokenHash.String()) == "" {
		return RefreshSession{}, domain.ErrInvalidRefreshTokenHash
	}
	if now.IsZero() {
		return RefreshSession{}, domain.ErrInvalidRefreshSessionTime
	}
	if ttl <= 0 {
		return RefreshSession{}, domain.ErrInvalidRefreshTokenTTL
	}

	return RefreshSession{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

func (s *RefreshSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *RefreshSession) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

func (s *RefreshSession) IsValid(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}

func (s *RefreshSession) EnsureUsable(now time.Time) error {
	if s.IsRevoked() {
		return domain.ErrRefreshTokenRevoked
	}
	if s.IsExpired(now) {
		return domain.ErrRefreshTokenExpired
	}
	return nil
}
