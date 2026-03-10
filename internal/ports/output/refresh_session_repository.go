package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type RefreshSessionRepository interface {
	Save(ctx context.Context, session entity.RefreshSession) error
	// FindByTokenHash используется для read-only lookup (logout, проверка существования).
	FindByTokenHash(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error)
	// FindByTokenHashForUpdate используется внутри транзакции при refresh rotation.
	// Выполняет SELECT ... FOR UPDATE чтобы исключить параллельный reuse одной сессии.
	FindByTokenHashForUpdate(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	DeleteExpiredAndRevoked(ctx context.Context, userID uuid.UUID) error
}
