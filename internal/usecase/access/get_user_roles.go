package access

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// GetUserRoles возвращает доменное представление ролей пользователя.
// Сначала проверяет существование пользователя, затем загружает его роли.
// Сборка доменного результата (dedup, sort, проверка битых связей)
// делегирована в domain/service.BuildUserRolesResult.
func (s *AccessService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domainservice.RoleView, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	userRoles, err := s.userRoleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: find user roles: %w", err)
	}

	if len(userRoles) == 0 {
		return []domainservice.RoleView{}, nil
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roles, err := s.roleRepo.FindByIDs(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("get user roles: find roles: %w", err)
	}

	return domainservice.BuildUserRolesResult(userRoles, roles)
}
