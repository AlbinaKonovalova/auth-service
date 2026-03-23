package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

var (
	testNow    = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	testUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

func validEmail(t *testing.T) value.Email {
	t.Helper()
	e, err := value.NewEmail("admin@example.com")
	require.NoError(t, err)
	return e
}

func TestNewUser_Valid(t *testing.T) {
	u, err := entity.NewUser(testUserID, validEmail(t), "hashed_password", testNow)
	require.NoError(t, err)
	assert.Equal(t, testUserID, u.ID)
	assert.Equal(t, "admin@example.com", u.Email)
	assert.True(t, u.IsActive)
}

func TestNewUser_NilID(t *testing.T) {
	_, err := entity.NewUser(uuid.Nil, validEmail(t), "hashed", testNow)
	assert.ErrorIs(t, err, domain.ErrInvalidUserID)
}

func TestNewUser_EmptyPasswordHash(t *testing.T) {
	_, err := entity.NewUser(testUserID, validEmail(t), "", testNow)
	assert.ErrorIs(t, err, domain.ErrInvalidPasswordHash)
}

func TestNewUser_WhitespacePasswordHash(t *testing.T) {
	_, err := entity.NewUser(testUserID, validEmail(t), "   ", testNow)
	assert.ErrorIs(t, err, domain.ErrInvalidPasswordHash)
}

func TestUser_EnsureActive_Active(t *testing.T) {
	u, _ := entity.NewUser(testUserID, validEmail(t), "hash", testNow)
	assert.NoError(t, u.EnsureActive())
}

func TestUser_EnsureActive_Inactive(t *testing.T) {
	u, _ := entity.NewUser(testUserID, validEmail(t), "hash", testNow)
	u.IsActive = false
	assert.ErrorIs(t, u.EnsureActive(), domain.ErrUserInactive)
}
