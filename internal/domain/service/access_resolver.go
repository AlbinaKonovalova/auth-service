package service

import (
	"sort"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func ResolveEffectiveAccess(
	userRoles []entity.UserRole,
	roles []entity.Role,
	rolePermissions []entity.RolePermission,
	permissions []entity.Permission,
) (roleCodes []string, permissionCodes []string, err error) {
	if len(userRoles) == 0 {
		return []string{}, []string{}, nil
	}

	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, role := range roles {
		roleByID[role.ID] = role
	}

	effectiveRoleIDs := make(map[uuid.UUID]struct{}, len(userRoles))
	effectiveRoleCodes := make(map[string]struct{}, len(userRoles))

	for _, userRole := range userRoles {
		role, ok := roleByID[userRole.RoleID]
		if !ok {
			return nil, nil, domain.ErrDataIntegrityViolation
		}

		effectiveRoleIDs[userRole.RoleID] = struct{}{}
		effectiveRoleCodes[role.Code] = struct{}{}
	}

	roleCodes = make([]string, 0, len(effectiveRoleCodes))
	for code := range effectiveRoleCodes {
		roleCodes = append(roleCodes, code)
	}
	sort.Strings(roleCodes)

	if len(rolePermissions) == 0 {
		return roleCodes, []string{}, nil
	}

	permissionByID := make(map[uuid.UUID]entity.Permission, len(permissions))
	for _, permission := range permissions {
		permissionByID[permission.ID] = permission
	}

	effectivePermissionCodes := make(map[string]struct{})

	for _, rolePermission := range rolePermissions {
		if _, ok := effectiveRoleIDs[rolePermission.RoleID]; !ok {
			continue
		}

		permission, ok := permissionByID[rolePermission.PermissionID]
		if !ok {
			return nil, nil, domain.ErrDataIntegrityViolation
		}

		effectivePermissionCodes[permission.Code] = struct{}{}
	}

	permissionCodes = make([]string, 0, len(effectivePermissionCodes))
	for code := range effectivePermissionCodes {
		permissionCodes = append(permissionCodes, code)
	}
	sort.Strings(permissionCodes)

	return roleCodes, permissionCodes, nil
}
