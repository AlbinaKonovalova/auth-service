package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// NewPasswordResetToken создаёт новый токен сброса пароля с доменной валидацией.
// tokenHash — уже захешированный токен (хешируется снаружи через token hasher).
// ttl задаёт срок жизни токена.
// Возвращает соответствующую доменную ошибку для каждой причины невалидности:
// - ErrInvalidResetTokenID если id == uuid.Nil
// - ErrInvalidResetTokenUserID если userID == uuid.Nil
// - ErrInvalidResetTokenHash если tokenHash пустой
// - ErrInvalidResetTokenNow если now нулевой
// - ErrInvalidResetTokenTTL если ttl <= 0
func NewPasswordResetToken(id, userID uuid.UUID, tokenHash string, now time.Time, ttl time.Duration) (PasswordResetToken, error) {
	if id == uuid.Nil {
		return PasswordResetToken{}, domain.ErrInvalidResetTokenID
	}
	if userID == uuid.Nil {
		return PasswordResetToken{}, domain.ErrInvalidResetTokenUserID
	}
	if strings.TrimSpace(tokenHash) == "" {
		return PasswordResetToken{}, domain.ErrInvalidResetTokenHash
	}
	if now.IsZero() {
		return PasswordResetToken{}, domain.ErrInvalidResetTokenNow
	}
	if ttl <= 0 {
		return PasswordResetToken{}, domain.ErrInvalidResetTokenTTL
	}

	return PasswordResetToken{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

func (t *PasswordResetToken) IsExpired(now time.Time) bool {
	return !now.Before(t.ExpiresAt)
}

func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

func (t *PasswordResetToken) EnsureUsable(now time.Time) error {
	if t.IsExpired(now) {
		return domain.ErrResetTokenExpired
	}
	if t.IsUsed() {
		return domain.ErrResetTokenUsed
	}
	return nil
}
