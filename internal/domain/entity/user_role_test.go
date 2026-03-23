package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func TestNewUserRole_Valid(t *testing.T) {
	ur, err := entity.NewUserRole(testUserID, uuid.New(), testNow)
	require.NoError(t, err)
	assert.Equal(t, testUserID, ur.UserID)
}

func TestNewUserRole_NilUserID(t *testing.T) {
	_, err := entity.NewUserRole(uuid.Nil, uuid.New(), testNow)
	assert.ErrorIs(t, err, domain.ErrInvalidUserRole)
}

func TestNewUserRole_NilRoleID(t *testing.T) {
	_, err := entity.NewUserRole(testUserID, uuid.Nil, testNow)
	assert.ErrorIs(t, err, domain.ErrInvalidUserRole)
}

func TestNewUserRole_ZeroTime(t *testing.T) {
	_, err := entity.NewUserRole(testUserID, uuid.New(), time.Time{})
	assert.ErrorIs(t, err, domain.ErrInvalidUserRole)
}
