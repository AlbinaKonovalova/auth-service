package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RoleRepository interface {
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Role, error)
	FindByCodes(ctx context.Context, codes []string) ([]entity.Role, error)
	FindByCode(ctx context.Context, code string) (*entity.Role, error)
	FindAll(ctx context.Context) ([]entity.Role, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Create(ctx context.Context, role entity.Role) error
}
