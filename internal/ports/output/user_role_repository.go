package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type UserRoleRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error)
}
