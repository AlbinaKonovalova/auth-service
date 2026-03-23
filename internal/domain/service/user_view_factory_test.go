package service_test

import (
	"testing"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func TestBuildAdminUserView_Valid(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	user := entity.User{ID: userID, Email: "admin@example.com", IsActive: true}
	roles := []entity.Role{{ID: roleID, Code: "admin", Name: "Admin"}}
	userRoles := []entity.UserRole{{UserID: userID, RoleID: roleID}}

	view, err := service.BuildAdminUserView(user, userRoles, roles)
	require.NoError(t, err)
	assert.Equal(t, userID, view.ID)
	assert.Equal(t, "admin@example.com", view.Email)
	assert.True(t, view.IsActive)
	assert.Equal(t, []string{"admin"}, view.Roles)
}

func TestBuildAdminUserView_MultipleRolesSorted(t *testing.T) {
	userID := uuid.New()
	roleID1, roleID2 := uuid.New(), uuid.New()
	user := entity.User{ID: userID, Email: "u@example.com", IsActive: true}
	roles := []entity.Role{
		{ID: roleID1, Code: "manager", Name: "Manager"},
		{ID: roleID2, Code: "admin", Name: "Admin"},
	}
	userRoles := []entity.UserRole{
		{UserID: userID, RoleID: roleID1},
		{UserID: userID, RoleID: roleID2},
	}

	view, err := service.BuildAdminUserView(user, userRoles, roles)
	require.NoError(t, err)
	assert.Equal(t, []string{"admin", "manager"}, view.Roles) // отсортированы
}

func TestBuildAdminUserView_BrokenRoleLink(t *testing.T) {
	userID := uuid.New()
	user := entity.User{ID: userID, Email: "u@example.com", IsActive: true}
	userRoles := []entity.UserRole{{UserID: userID, RoleID: uuid.New()}} // несуществующая роль

	_, err := service.BuildAdminUserView(user, userRoles, nil)
	assert.ErrorIs(t, err, domain.ErrDataIntegrityViolation)
}

func TestBuildUserList_Empty(t *testing.T) {
	list, err := service.BuildUserList(nil, nil, nil, 0, 1, 20)
	require.NoError(t, err)
	assert.Empty(t, list.Items)
	assert.Equal(t, 0, list.Total)
}

func TestBuildUserList_NegativeTotal(t *testing.T) {
	_, err := service.BuildUserList(nil, nil, nil, -1, 1, 20)
	assert.ErrorIs(t, err, domain.ErrInvalidUserList)
}

func TestNormalizeUserListPagination_Defaults(t *testing.T) {
	page, perPage := service.NormalizeUserListPagination(0, 0)
	assert.Equal(t, service.DefaultUserListPage, page)
	assert.Equal(t, service.DefaultUserListPerPage, perPage)
}

func TestNormalizeUserListPagination_ExceedsMax(t *testing.T) {
	_, perPage := service.NormalizeUserListPagination(1, 999)
	assert.Equal(t, service.MaxUserListPerPage, perPage)
}

func TestNormalizeUserListPagination_Valid(t *testing.T) {
	page, perPage := service.NormalizeUserListPagination(3, 50)
	assert.Equal(t, 3, page)
	assert.Equal(t, 50, perPage)
}

func TestValidateUserListFilters_NoRole(t *testing.T) {
	filters := dto.UserListFilters{Page: 0, PerPage: 0}
	result, err := service.ValidateUserListFilters(filters, nil)
	require.NoError(t, err)
	assert.Equal(t, service.DefaultUserListPage, result.Page)
	assert.Equal(t, service.DefaultUserListPerPage, result.PerPage)
}

func TestValidateUserListFilters_ValidRole(t *testing.T) {
	filters := dto.UserListFilters{Role: "ADMIN", Page: 1, PerPage: 10}
	foundRoles := []entity.Role{{Code: "admin"}}
	result, err := service.ValidateUserListFilters(filters, foundRoles)
	require.NoError(t, err)
	assert.Equal(t, "admin", result.Role) // нормализован
}

func TestValidateUserListFilters_RoleNotFound(t *testing.T) {
	filters := dto.UserListFilters{Role: "admin"}
	_, err := service.ValidateUserListFilters(filters, nil) // роль не найдена в БД
	assert.ErrorIs(t, err, domain.ErrRoleNotFound)
}

func TestValidateUserListFilters_InvalidRoleCode(t *testing.T) {
	filters := dto.UserListFilters{Role: ""}
	// пустой role — игнорируется (фильтр по роли не задан)
	result, err := service.ValidateUserListFilters(filters, nil)
	require.NoError(t, err)
	assert.Equal(t, "", result.Role)
}
