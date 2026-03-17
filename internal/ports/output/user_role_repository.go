package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type UserRoleRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error)
	FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]entity.UserRole, error)
	Assign(ctx context.Context, userRole entity.UserRole) error
}
