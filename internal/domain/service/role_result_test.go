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

func TestValidateRoleDeletion_NotInUse(t *testing.T) {
	assert.NoError(t, service.ValidateRoleDeletion(false, false))
}

func TestValidateRoleDeletion_HasUsers(t *testing.T) {
	assert.ErrorIs(t, service.ValidateRoleDeletion(true, false), domain.ErrRoleInUse)
}

func TestValidateRoleDeletion_HasPermissions(t *testing.T) {
	assert.ErrorIs(t, service.ValidateRoleDeletion(false, true), domain.ErrRoleInUse)
}

func TestValidateRoleDeletion_BothInUse(t *testing.T) {
	assert.ErrorIs(t, service.ValidateRoleDeletion(true, true), domain.ErrRoleInUse)
}

func TestBuildRoleListResult_SortedByCode(t *testing.T) {
	roles := []entity.Role{
		{ID: uuid.New(), Code: "manager", Name: "Manager"},
		{ID: uuid.New(), Code: "admin", Name: "Admin"},
		{ID: uuid.New(), Code: "analyst", Name: "Analyst"},
	}
	result := service.BuildRoleListResult(roles)
	require.Len(t, result, 3)
	assert.Equal(t, "admin", result[0].Code)
	assert.Equal(t, "analyst", result[1].Code)
	assert.Equal(t, "manager", result[2].Code)
}

func TestBuildUserRolesResult_Valid(t *testing.T) {
	roleID := uuid.New()
	roles := []entity.Role{{ID: roleID, Code: "admin", Name: "Admin"}}
	userRoles := []entity.UserRole{{UserID: uuid.New(), RoleID: roleID}}

	result, err := service.BuildUserRolesResult(userRoles, roles)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "admin", result[0].Code)
}

func TestBuildUserRolesResult_Deduplicated(t *testing.T) {
	roleID := uuid.New()
	roles := []entity.Role{{ID: roleID, Code: "admin", Name: "Admin"}}
	userRoles := []entity.UserRole{
		{UserID: uuid.New(), RoleID: roleID},
		{UserID: uuid.New(), RoleID: roleID}, // дубль
	}

	result, err := service.BuildUserRolesResult(userRoles, roles)
	require.NoError(t, err)
	assert.Len(t, result, 1) // дубль убран
}

func TestBuildUserRolesResult_BrokenLink(t *testing.T) {
	// role_id есть в user_roles, но нет в справочнике — data integrity violation
	userRoles := []entity.UserRole{{UserID: uuid.New(), RoleID: uuid.New()}}
	_, err := service.BuildUserRolesResult(userRoles, nil)
	assert.ErrorIs(t, err, domain.ErrDataIntegrityViolation)
}

func TestBuildUserRolesResult_Empty(t *testing.T) {
	result, err := service.BuildUserRolesResult(nil, nil)
	require.NoError(t, err)
	assert.Empty(t, result)
}
