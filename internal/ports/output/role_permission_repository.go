package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RolePermissionRepository interface {
	FindByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]entity.RolePermission, error)
}
