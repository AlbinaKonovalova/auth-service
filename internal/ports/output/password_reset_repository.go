package output

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, token entity.PasswordResetToken) error

	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)

	FindByTokenHashForUpdate(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)

	InvalidateByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error

	MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error
}
