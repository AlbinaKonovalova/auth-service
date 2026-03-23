package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func TestNewPermission_Valid(t *testing.T) {
	p, err := entity.NewPermission(entity.NewPermissionParams{
		ID: uuid.New(), Code: "users.read", Description: "Read users",
	})
	require.NoError(t, err)
	assert.Equal(t, "Read users", p.Description)
}

func TestNewPermission_EmptyDescription(t *testing.T) {
	_, err := entity.NewPermission(entity.NewPermissionParams{
		ID: uuid.New(), Code: "users.read", Description: "",
	})
	assert.ErrorIs(t, err, domain.ErrPermissionDescriptionEmpty)
}

func TestNewPermission_WhitespaceDescription(t *testing.T) {
	_, err := entity.NewPermission(entity.NewPermissionParams{
		ID: uuid.New(), Code: "users.read", Description: "   ",
	})
	assert.ErrorIs(t, err, domain.ErrPermissionDescriptionEmpty)
}
