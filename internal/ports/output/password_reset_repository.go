package output

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

// PasswordResetRepository — контракт для хранилища токенов сброса пароля.
type PasswordResetRepository interface {
	// Create сохраняет новый токен сброса пароля.
	Create(ctx context.Context, token entity.PasswordResetToken) error

	// FindByTokenHash ищет токен по хешу (без блокировки, для read-only сценариев).
	// Возвращает domain.ErrResetTokenNotFound если токен не найден.
	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)

	// FindByTokenHashForUpdate ищет токен по хешу с блокировкой строки (FOR UPDATE).
	// Используется в confirm-сценарии для сериализации параллельных confirm по одному токену.
	// Возвращает domain.ErrResetTokenNotFound если токен не найден.
	FindByTokenHashForUpdate(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)

	// InvalidateByUserID аннулирует все незавершённые токены пользователя,
	// записывая им UsedAt = now. В модели сервиса нет отдельного состояния "revoked" —
	// аннулирование при создании нового токена кодируется через то же поле UsedAt.
	// Вызывается перед созданием нового reset token, чтобы старые запросы перестали работать.
	InvalidateByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error

	// MarkUsed помечает токен как использованный, записывая время использования.
	// Возвращает domain.ErrResetTokenUsed если токен уже был помечен (параллельный confirm).
	MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error
}
