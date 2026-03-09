package entity

import (
	"time"

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

func (s *RefreshSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *RefreshSession) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

func (s *RefreshSession) IsValid(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}
