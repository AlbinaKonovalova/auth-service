package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

// PasswordResetToken представляет одноразовый токен для сброса пароля.
// Токен считается применимым (usable), если он не истёк и не был использован ранее.
// В модели сервиса нет отдельного состояния "revoked" —
// аннулирование предыдущих токенов при новом запросе кодируется через UsedAt (invalidation as used).
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

// IsExpired возвращает true если токен истёк к моменту now.
func (t *PasswordResetToken) IsExpired(now time.Time) bool {
	return !now.Before(t.ExpiresAt)
}

// IsUsed возвращает true если токен уже был использован или аннулирован.
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// EnsureUsable проверяет, что токен можно применить.
// Возвращает доменную ошибку если токен истёк или уже использован.
// Порядок проверок зафиксирован: сначала expired, затем used —
// истёкший токен не должен раскрывать информацию о том, был ли он использован.
func (t *PasswordResetToken) EnsureUsable(now time.Time) error {
	if t.IsExpired(now) {
		return domain.ErrResetTokenExpired
	}
	if t.IsUsed() {
		return domain.ErrResetTokenUsed
	}
	return nil
}
