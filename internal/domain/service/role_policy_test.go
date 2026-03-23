package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func TestCanRevokeUserRole_Valid(t *testing.T) {
	roleID1, roleID2 := uuid.New(), uuid.New()
	userRoles := []entity.UserRole{{RoleID: roleID1}, {RoleID: roleID2}}
	assert.NoError(t, service.CanRevokeUserRole(userRoles, roleID1))
}

func TestCanRevokeUserRole_NotAssigned(t *testing.T) {
	roleID1, roleID2 := uuid.New(), uuid.New()
	userRoles := []entity.UserRole{{RoleID: roleID1}}
	assert.ErrorIs(t, service.CanRevokeUserRole(userRoles, roleID2), domain.ErrUserRoleNotFound)
}

func TestCanRevokeUserRole_LastRole(t *testing.T) {
	roleID := uuid.New()
	userRoles := []entity.UserRole{{RoleID: roleID}}
	assert.ErrorIs(t, service.CanRevokeUserRole(userRoles, roleID), domain.ErrCannotRevokeLastRole)
}

func TestCanRevokeUserRole_EmptyRoles(t *testing.T) {
	assert.ErrorIs(t, service.CanRevokeUserRole(nil, uuid.New()), domain.ErrUserRoleNotFound)
}
