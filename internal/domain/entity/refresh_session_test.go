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

func TestNewRefreshSession_Valid(t *testing.T) {
	s, err := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("somehash"), testNow, time.Hour)
	require.NoError(t, err)
	assert.Equal(t, testNow.Add(time.Hour), s.ExpiresAt)
	assert.Nil(t, s.RevokedAt)
}

func TestNewRefreshSession_NilSessionID(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.Nil, testUserID, value.TokenHash("hash"), testNow, time.Hour)
	assert.ErrorIs(t, err, domain.ErrInvalidRefreshSessionID)
}

func TestNewRefreshSession_NilUserID(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.New(), uuid.Nil, value.TokenHash("hash"), testNow, time.Hour)
	assert.ErrorIs(t, err, domain.ErrInvalidUserID)
}

func TestNewRefreshSession_EmptyTokenHash(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash(""), testNow, time.Hour)
	assert.ErrorIs(t, err, domain.ErrInvalidRefreshTokenHash)
}

func TestNewRefreshSession_ZeroTime(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("hash"), time.Time{}, time.Hour)
	assert.ErrorIs(t, err, domain.ErrInvalidRefreshSessionTime)
}

func TestNewRefreshSession_ZeroTTL(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("hash"), testNow, 0)
	assert.ErrorIs(t, err, domain.ErrInvalidRefreshTokenTTL)
}

func TestNewRefreshSession_NegativeTTL(t *testing.T) {
	_, err := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("hash"), testNow, -time.Hour)
	assert.ErrorIs(t, err, domain.ErrInvalidRefreshTokenTTL)
}

func TestRefreshSession_IsRevoked(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	assert.False(t, s.IsRevoked())

	revokedAt := testNow
	s.RevokedAt = &revokedAt
	assert.True(t, s.IsRevoked())
}

func TestRefreshSession_IsExpired(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	assert.False(t, s.IsExpired(testNow))
	assert.False(t, s.IsExpired(testNow.Add(59*time.Minute)))
	assert.True(t, s.IsExpired(testNow.Add(time.Hour)))
	assert.True(t, s.IsExpired(testNow.Add(2*time.Hour)))
}

func TestRefreshSession_EnsureUsable_Valid(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	assert.NoError(t, s.EnsureUsable(testNow.Add(30*time.Minute)))
}

func TestRefreshSession_EnsureUsable_Revoked(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	revokedAt := testNow
	s.RevokedAt = &revokedAt
	assert.ErrorIs(t, s.EnsureUsable(testNow), domain.ErrRefreshTokenRevoked)
}

func TestRefreshSession_EnsureUsable_Expired(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	assert.ErrorIs(t, s.EnsureUsable(testNow.Add(2*time.Hour)), domain.ErrRefreshTokenExpired)
}

func TestRefreshSession_EnsureUsable_RevokedAndExpired_RevokedWins(t *testing.T) {
	s, _ := entity.NewRefreshSession(uuid.New(), testUserID, value.TokenHash("h"), testNow, time.Hour)
	revokedAt := testNow
	s.RevokedAt = &revokedAt
	assert.ErrorIs(t, s.EnsureUsable(testNow.Add(2*time.Hour)), domain.ErrRefreshTokenRevoked)
}
