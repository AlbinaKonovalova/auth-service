package common

import (
	"context"

	"github.com/google/uuid"

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

func (r *PermissionResolver) Resolve(ctx context.Context, userID uuid.UUID) (roles []string, perms []string, err error) {
	userRoles, err := r.userRoles.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if len(userRoles) == 0 {
		return nil, nil, nil
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roleEntities, err := r.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return nil, nil, err
	}
	for _, role := range roleEntities {
		roles = append(roles, role.Code)
	}

	rolePerms, err := r.rolePermissions.FindByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, nil, err
	}
	if len(rolePerms) == 0 {
		return roles, nil, nil
	}

	seen := make(map[uuid.UUID]struct{}, len(rolePerms))
	permIDs := make([]uuid.UUID, 0, len(rolePerms))
	for _, rp := range rolePerms {
		if _, ok := seen[rp.PermissionID]; !ok {
			seen[rp.PermissionID] = struct{}{}
			permIDs = append(permIDs, rp.PermissionID)
		}
	}

	permEntities, err := r.permissions.FindByIDs(ctx, permIDs)
	if err != nil {
		return nil, nil, err
	}
	for _, p := range permEntities {
		perms = append(perms, p.Code)
	}

	return roles, perms, nil
}
