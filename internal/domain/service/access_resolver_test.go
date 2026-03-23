package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func TestResolveEffectiveAccess_NoRoles(t *testing.T) {
	roles, perms, err := service.ResolveEffectiveAccess(nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Empty(t, roles)
	assert.Empty(t, perms)
}

func TestResolveEffectiveAccess_RolesNoPermissions(t *testing.T) {
	roleID := uuid.New()
	userRoles := []entity.UserRole{{RoleID: roleID}}
	roles := []entity.Role{{ID: roleID, Code: "admin"}}

	roleCodes, permCodes, err := service.ResolveEffectiveAccess(userRoles, roles, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"admin"}, roleCodes)
	assert.Empty(t, permCodes)
}

func TestResolveEffectiveAccess_WithPermissions(t *testing.T) {
	roleID := uuid.New()
	permID1, permID2 := uuid.New(), uuid.New()

	userRoles := []entity.UserRole{{RoleID: roleID}}
	roles := []entity.Role{{ID: roleID, Code: "admin"}}
	rolePerms := []entity.RolePermission{
		{RoleID: roleID, PermissionID: permID1},
		{RoleID: roleID, PermissionID: permID2},
	}
	perms := []entity.Permission{
		{ID: permID1, Code: "users.write"},
		{ID: permID2, Code: "users.read"},
	}

	roleCodes, permCodes, err := service.ResolveEffectiveAccess(userRoles, roles, rolePerms, perms)
	require.NoError(t, err)
	assert.Equal(t, []string{"admin"}, roleCodes)
	assert.Equal(t, []string{"users.read", "users.write"}, permCodes) // отсортированы
}

func TestResolveEffectiveAccess_DedupPermissions(t *testing.T) {
	roleID1, roleID2 := uuid.New(), uuid.New()
	permID := uuid.New()

	userRoles := []entity.UserRole{{RoleID: roleID1}, {RoleID: roleID2}}
	roles := []entity.Role{
		{ID: roleID1, Code: "admin"},
		{ID: roleID2, Code: "manager"},
	}
	rolePerms := []entity.RolePermission{
		{RoleID: roleID1, PermissionID: permID},
		{RoleID: roleID2, PermissionID: permID}, // одно и то же permission в двух ролях
	}
	perms := []entity.Permission{{ID: permID, Code: "users.read"}}

	_, permCodes, err := service.ResolveEffectiveAccess(userRoles, roles, rolePerms, perms)
	require.NoError(t, err)
	assert.Equal(t, []string{"users.read"}, permCodes) // дубль убран
}

func TestResolveEffectiveAccess_BrokenRoleLink(t *testing.T) {
	// role_id есть в user_roles, но нет в справочнике roles
	userRoles := []entity.UserRole{{RoleID: uuid.New()}}
	_, _, err := service.ResolveEffectiveAccess(userRoles, nil, nil, nil)
	assert.ErrorIs(t, err, domain.ErrDataIntegrityViolation)
}

func TestResolveEffectiveAccess_BrokenPermissionLink(t *testing.T) {
	roleID := uuid.New()
	permID := uuid.New()

	userRoles := []entity.UserRole{{RoleID: roleID}}
	roles := []entity.Role{{ID: roleID, Code: "admin"}}
	rolePerms := []entity.RolePermission{{RoleID: roleID, PermissionID: permID}}
	// permission отсутствует в справочнике
	_, _, err := service.ResolveEffectiveAccess(userRoles, roles, rolePerms, nil)
	assert.ErrorIs(t, err, domain.ErrDataIntegrityViolation)
}

func TestResolveEffectiveAccess_PermissionOfUnrelatedRole_Ignored(t *testing.T) {
	// у пользователя только roleID1, но rolePerms содержит связь с roleID2
	roleID1 := uuid.New()
	roleID2 := uuid.New()
	permID := uuid.New()

	userRoles := []entity.UserRole{{RoleID: roleID1}}
	roles := []entity.Role{{ID: roleID1, Code: "admin"}}
	rolePerms := []entity.RolePermission{
		{RoleID: roleID2, PermissionID: permID}, // другая роль — должна игнорироваться
	}
	perms := []entity.Permission{{ID: permID, Code: "users.read"}}

	_, permCodes, err := service.ResolveEffectiveAccess(userRoles, roles, rolePerms, perms)
	require.NoError(t, err)
	assert.Empty(t, permCodes) // permission другой роли не попал
}
