package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestNormalizeRequestedRoleCodes_Valid(t *testing.T) {
	codes, err := entity.NormalizeRequestedRoleCodes([]string{"Admin", "manager"})
	require.NoError(t, err)
	require.Len(t, codes, 2)
	assert.Equal(t, "admin", codes[0].String())
	assert.Equal(t, "manager", codes[1].String())
}

func TestNormalizeRequestedRoleCodes_Empty(t *testing.T) {
	codes, err := entity.NormalizeRequestedRoleCodes(nil)
	require.NoError(t, err)
	assert.Nil(t, codes)
}

func TestNormalizeRequestedRoleCodes_Duplicate(t *testing.T) {
	_, err := entity.NormalizeRequestedRoleCodes([]string{"admin", "ADMIN"})
	assert.ErrorIs(t, err, domain.ErrDuplicateRoleCode)
}

func TestNormalizeRequestedRoleCodes_InvalidCode(t *testing.T) {
	_, err := entity.NormalizeRequestedRoleCodes([]string{"admin", ""})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestEnsureAllRequestedRolesExist_AllFound(t *testing.T) {
	adminCode, _ := value.NewRoleCode("admin")
	requested := []value.RoleCode{adminCode}
	found := []entity.Role{{ID: uuid.New(), Code: "admin", Name: "Admin"}}
	assert.NoError(t, entity.EnsureAllRequestedRolesExist(requested, found))
}

func TestEnsureAllRequestedRolesExist_NotFound(t *testing.T) {
	adminCode, _ := value.NewRoleCode("admin")
	requested := []value.RoleCode{adminCode}
	found := []entity.Role{{ID: uuid.New(), Code: "manager", Name: "Manager"}}
	assert.ErrorIs(t, entity.EnsureAllRequestedRolesExist(requested, found), domain.ErrRoleNotFound)
}

func TestEnsureAllRequestedRolesExist_EmptyRequested(t *testing.T) {
	assert.NoError(t, entity.EnsureAllRequestedRolesExist(nil, nil))
}
