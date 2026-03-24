package service

import (
	"sort"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type PermissionView struct {
	ID          uuid.UUID
	Code        string
	Description string
}

func PermissionViewFromEntity(p entity.Permission) PermissionView {
	return PermissionView{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
	}
}

func ValidatePermissionDeletion(hasRoles bool) error {
	if hasRoles {
		return domain.ErrPermissionInUse
	}
	return nil
}

func BuildPermissionListResult(permissions []entity.Permission) []PermissionView {
	result := make([]PermissionView, len(permissions))
	for i, p := range permissions {
		result[i] = PermissionViewFromEntity(p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result
}

func BuildRolePermissionsResult(rolePerms []entity.RolePermission, permissions []entity.Permission) ([]PermissionView, error) {
	if len(rolePerms) == 0 {
		return []PermissionView{}, nil
	}

	permByID := make(map[uuid.UUID]entity.Permission, len(permissions))
	for _, p := range permissions {
		permByID[p.ID] = p
	}

	seen := make(map[uuid.UUID]struct{}, len(rolePerms))
	result := make([]PermissionView, 0, len(rolePerms))

	for _, rp := range rolePerms {
		if _, dup := seen[rp.PermissionID]; dup {
			continue
		}
		seen[rp.PermissionID] = struct{}{}

		p, ok := permByID[rp.PermissionID]
		if !ok {
			return nil, domain.ErrDataIntegrityViolation
		}

		result = append(result, PermissionViewFromEntity(p))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result, nil
}
