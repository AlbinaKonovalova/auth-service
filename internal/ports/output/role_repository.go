package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RoleRepository interface {
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Role, error)
	FindByCodes(ctx context.Context, codes []string) ([]entity.Role, error)
}
