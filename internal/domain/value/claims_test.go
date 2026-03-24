package value_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestNewAccessClaims_Valid(t *testing.T) {
	userID := uuid.New()
	claims, err := value.NewAccessClaims(userID, "admin@example.com", []string{"admin"}, []string{"users.read"})
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "admin@example.com", claims.Email)
	assert.Equal(t, []string{"admin"}, claims.Roles)
	assert.Equal(t, []string{"users.read"}, claims.Permissions)
}

func TestNewAccessClaims_NilUserID(t *testing.T) {
	_, err := value.NewAccessClaims(uuid.Nil, "admin@example.com", nil, nil)
	assert.ErrorIs(t, err, domain.ErrInvalidUserID)
}

func TestNewAccessClaims_EmptyEmail(t *testing.T) {
	_, err := value.NewAccessClaims(uuid.New(), "", nil, nil)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestNewAccessClaims_WhitespaceEmail(t *testing.T) {
	_, err := value.NewAccessClaims(uuid.New(), "   ", nil, nil)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestNewAccessClaims_NilRolesAndPermissions_Valid(t *testing.T) {
	claims, err := value.NewAccessClaims(uuid.New(), "user@example.com", nil, nil)
	require.NoError(t, err)
	assert.Nil(t, claims.Roles)
	assert.Nil(t, claims.Permissions)
}
