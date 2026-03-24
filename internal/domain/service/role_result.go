package service

import (
	"sort"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

type RoleView struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

func RoleViewFromEntity(r entity.Role) RoleView {
	return RoleView{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
	}
}

func ValidateRoleDeletion(hasUsers bool, hasPermissions bool) error {
	if hasUsers || hasPermissions {
		return domain.ErrRoleInUse
	}
	return nil
}

func BuildRoleListResult(roles []entity.Role) []RoleView {
	result := make([]RoleView, len(roles))
	for i, r := range roles {
		result[i] = RoleViewFromEntity(r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result
}

func BuildUserRolesResult(userRoles []entity.UserRole, roles []entity.Role) ([]RoleView, error) {
	if len(userRoles) == 0 {
		return []RoleView{}, nil
	}

	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, r := range roles {
		roleByID[r.ID] = r
	}

	seen := make(map[uuid.UUID]struct{}, len(userRoles))
	result := make([]RoleView, 0, len(userRoles))

	for _, ur := range userRoles {
		if _, dup := seen[ur.RoleID]; dup {
			continue
		}
		seen[ur.RoleID] = struct{}{}

		r, ok := roleByID[ur.RoleID]
		if !ok {
			return nil, domain.ErrDataIntegrityViolation
		}

		result = append(result, RoleViewFromEntity(r))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result, nil
}
