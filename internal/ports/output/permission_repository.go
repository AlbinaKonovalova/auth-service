package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type PermissionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Permission, error)
	FindByCode(ctx context.Context, code string) (*entity.Permission, error)
	FindAll(ctx context.Context) ([]entity.Permission, error)
}
