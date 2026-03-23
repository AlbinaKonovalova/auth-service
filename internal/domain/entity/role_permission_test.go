package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func TestNewRolePermission_Valid(t *testing.T) {
	rp, err := entity.NewRolePermission(uuid.New(), uuid.New())
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, rp.RoleID)
	assert.NotEqual(t, uuid.Nil, rp.PermissionID)
}

func TestNewRolePermission_NilRoleID(t *testing.T) {
	_, err := entity.NewRolePermission(uuid.Nil, uuid.New())
	assert.ErrorIs(t, err, domain.ErrInvalidRolePermission)
}

func TestNewRolePermission_NilPermissionID(t *testing.T) {
	_, err := entity.NewRolePermission(uuid.New(), uuid.Nil)
	assert.ErrorIs(t, err, domain.ErrInvalidRolePermission)
}
