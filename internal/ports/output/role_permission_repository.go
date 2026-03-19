package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RolePermissionRepository interface {
	FindByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.RolePermission, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]entity.RolePermission, error)
	Assign(ctx context.Context, rp entity.RolePermission) error
	Exists(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)
	Revoke(ctx context.Context, roleID, permissionID uuid.UUID) error
}
