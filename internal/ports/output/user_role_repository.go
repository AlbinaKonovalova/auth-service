package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type UserRoleRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error)
	FindByUserIDForUpdate(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error)
	FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]entity.UserRole, error)
	Exists(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	Assign(ctx context.Context, userRole entity.UserRole) error
	Revoke(ctx context.Context, userID, roleID uuid.UUID) error
}
