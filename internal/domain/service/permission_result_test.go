package service_test

import (
	"testing"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePermissionDeletion_NotInUse(t *testing.T) {
	assert.NoError(t, service.ValidatePermissionDeletion(false))
}

func TestValidatePermissionDeletion_InUse(t *testing.T) {
	assert.ErrorIs(t, service.ValidatePermissionDeletion(true), domain.ErrPermissionInUse)
}

func TestBuildPermissionListResult_SortedByCode(t *testing.T) {
	perms := []entity.Permission{
		{ID: uuid.New(), Code: "users.write", Description: "Write"},
		{ID: uuid.New(), Code: "stats.export", Description: "Export"},
		{ID: uuid.New(), Code: "users.read", Description: "Read"},
	}
	result := service.BuildPermissionListResult(perms)
	require.Len(t, result, 3)
	assert.Equal(t, "stats.export", result[0].Code)
	assert.Equal(t, "users.read", result[1].Code)
	assert.Equal(t, "users.write", result[2].Code)
}

func TestBuildRolePermissionsResult_Valid(t *testing.T) {
	permID := uuid.New()
	roleID := uuid.New()
	rolePerms := []entity.RolePermission{{RoleID: roleID, PermissionID: permID}}
	perms := []entity.Permission{{ID: permID, Code: "users.read", Description: "Read"}}

	result, err := service.BuildRolePermissionsResult(rolePerms, perms)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "users.read", result[0].Code)
}

func TestBuildRolePermissionsResult_Deduplicated(t *testing.T) {
	permID := uuid.New()
	roleID := uuid.New()
	rolePerms := []entity.RolePermission{
		{RoleID: roleID, PermissionID: permID},
		{RoleID: roleID, PermissionID: permID}, // дубль
	}
	perms := []entity.Permission{{ID: permID, Code: "users.read", Description: "Read"}}

	result, err := service.BuildRolePermissionsResult(rolePerms, perms)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestBuildRolePermissionsResult_BrokenLink(t *testing.T) {
	rolePerms := []entity.RolePermission{{RoleID: uuid.New(), PermissionID: uuid.New()}}
	_, err := service.BuildRolePermissionsResult(rolePerms, nil)
	assert.ErrorIs(t, err, domain.ErrDataIntegrityViolation)
}

func TestBuildRolePermissionsResult_Empty(t *testing.T) {
	result, err := service.BuildRolePermissionsResult(nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result)
}
