package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func TestNewRole_Valid(t *testing.T) {
	r, err := entity.NewRole(entity.NewRoleParams{
		ID: uuid.New(), Code: "admin", Name: "Admin", Description: "Full access",
	})
	require.NoError(t, err)
	assert.Equal(t, "Admin", r.Name)
	assert.Equal(t, "admin", r.Code)
}

func TestNewRole_EmptyName(t *testing.T) {
	_, err := entity.NewRole(entity.NewRoleParams{
		ID: uuid.New(), Code: "admin", Name: "", Description: "x",
	})
	assert.ErrorIs(t, err, domain.ErrRoleNameEmpty)
}

func TestNewRole_WhitespaceName(t *testing.T) {
	_, err := entity.NewRole(entity.NewRoleParams{
		ID: uuid.New(), Code: "admin", Name: "   ", Description: "x",
	})
	assert.ErrorIs(t, err, domain.ErrRoleNameEmpty)
}
