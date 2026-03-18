package common

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type PermissionResolver struct {
	userRoles       output.UserRoleRepository
	rolePermissions output.RolePermissionRepository
	roles           output.RoleRepository
	permissions     output.PermissionRepository
}

func NewPermissionResolver(
	userRoles output.UserRoleRepository,
	rolePermissions output.RolePermissionRepository,
	roles output.RoleRepository,
	permissions output.PermissionRepository,
) *PermissionResolver {
	return &PermissionResolver{
		userRoles:       userRoles,
		rolePermissions: rolePermissions,
		roles:           roles,
		permissions:     permissions,
	}
}

func (r *PermissionResolver) Resolve(ctx context.Context, userID uuid.UUID) ([]string, []string, error) {
	userRoleRows, err := r.userRoles.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if len(userRoleRows) == 0 {
		return []string{}, []string{}, nil
	}

	roleIDs := uniqueRoleIDs(userRoleRows)

	roleEntities, err := r.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return nil, nil, err
	}

	rolePermissionRows, err := r.rolePermissions.FindByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, nil, err
	}

	permissionEntities := []entity.Permission{}
	permissionIDs := uniquePermissionIDs(rolePermissionRows)
	if len(permissionIDs) > 0 {
		permissionEntities, err = r.permissions.FindByIDs(ctx, permissionIDs)
		if err != nil {
			return nil, nil, err
		}
	}

	return domainservice.ResolveEffectiveAccess(
		userRoleRows,
		roleEntities,
		rolePermissionRows,
		permissionEntities,
	)
}

func uniqueRoleIDs(userRoles []entity.UserRole) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(userRoles))
	seen := make(map[uuid.UUID]struct{}, len(userRoles))

	for _, userRole := range userRoles {
		if _, ok := seen[userRole.RoleID]; ok {
			continue
		}
		seen[userRole.RoleID] = struct{}{}
		result = append(result, userRole.RoleID)
	}

	return result
}

func uniquePermissionIDs(rolePermissions []entity.RolePermission) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(rolePermissions))
	seen := make(map[uuid.UUID]struct{}, len(rolePermissions))

	for _, rolePermission := range rolePermissions {
		if _, ok := seen[rolePermission.PermissionID]; ok {
			continue
		}
		seen[rolePermission.PermissionID] = struct{}{}
		result = append(result, rolePermission.PermissionID)
	}

	return result
}
