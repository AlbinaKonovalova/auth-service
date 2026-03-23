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

func TestNewPasswordResetToken_Valid(t *testing.T) {
	tokenID := uuid.New()
	tok, err := entity.NewPasswordResetToken(tokenID, testUserID, "hashedtoken", testNow, 30*time.Minute)
	require.NoError(t, err)
	assert.Equal(t, tokenID, tok.ID)
	assert.Equal(t, testNow.Add(30*time.Minute), tok.ExpiresAt)
	assert.Nil(t, tok.UsedAt)
}

func TestNewPasswordResetToken_NilID(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.Nil, testUserID, "hash", testNow, time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenID)
}

func TestNewPasswordResetToken_NilUserID(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), uuid.Nil, "hash", testNow, time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenUserID)
}

func TestNewPasswordResetToken_EmptyHash(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), testUserID, "", testNow, time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenHash)
}

func TestNewPasswordResetToken_WhitespaceHash(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), testUserID, "   ", testNow, time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenHash)
}

func TestNewPasswordResetToken_ZeroTime(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", time.Time{}, time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenNow)
}

func TestNewPasswordResetToken_ZeroTTL(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 0)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenTTL)
}

func TestNewPasswordResetToken_NegativeTTL(t *testing.T) {
	_, err := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, -time.Minute)
	assert.ErrorIs(t, err, domain.ErrInvalidResetTokenTTL)
}

func TestPasswordResetToken_IsExpired(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	assert.False(t, tok.IsExpired(testNow))
	assert.False(t, tok.IsExpired(testNow.Add(29*time.Minute)))
	assert.True(t, tok.IsExpired(testNow.Add(30*time.Minute)))
	assert.True(t, tok.IsExpired(testNow.Add(time.Hour)))
}

func TestPasswordResetToken_IsUsed(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	assert.False(t, tok.IsUsed())

	usedAt := testNow
	tok.UsedAt = &usedAt
	assert.True(t, tok.IsUsed())
}

func TestPasswordResetToken_EnsureUsable_Valid(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	assert.NoError(t, tok.EnsureUsable(testNow.Add(5*time.Minute)))
}

func TestPasswordResetToken_EnsureUsable_Expired(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	assert.ErrorIs(t, tok.EnsureUsable(testNow.Add(time.Hour)), domain.ErrResetTokenExpired)
}

func TestPasswordResetToken_EnsureUsable_Used(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	usedAt := testNow
	tok.UsedAt = &usedAt
	assert.ErrorIs(t, tok.EnsureUsable(testNow.Add(5*time.Minute)), domain.ErrResetTokenUsed)
}

// Expired проверяется раньше Used — если token и истёк, и использован, побеждает Expired.
func TestPasswordResetToken_EnsureUsable_ExpiredBeforeUsed(t *testing.T) {
	tok, _ := entity.NewPasswordResetToken(uuid.New(), testUserID, "hash", testNow, 30*time.Minute)
	usedAt := testNow.Add(5 * time.Minute)
	tok.UsedAt = &usedAt
	assert.ErrorIs(t, tok.EnsureUsable(testNow.Add(time.Hour)), domain.ErrResetTokenExpired)
}
