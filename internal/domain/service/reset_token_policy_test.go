package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func makeResetToken(now time.Time, ttl time.Duration, usedAt *time.Time) *entity.PasswordResetToken {
	return &entity.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: "hash",
		ExpiresAt: now.Add(ttl),
		UsedAt:    usedAt,
		CreatedAt: now,
	}
}

func TestValidateResetToken_Valid(t *testing.T) {
	tok := makeResetToken(testNow, 30*time.Minute, nil)
	assert.NoError(t, service.ValidateResetToken(tok, testNow.Add(5*time.Minute)))
}

func TestValidateResetToken_Expired(t *testing.T) {
	tok := makeResetToken(testNow, 30*time.Minute, nil)
	assert.ErrorIs(t, service.ValidateResetToken(tok, testNow.Add(time.Hour)), domain.ErrResetTokenExpired)
}

func TestValidateResetToken_Used(t *testing.T) {
	usedAt := testNow
	tok := makeResetToken(testNow, 30*time.Minute, &usedAt)
	assert.ErrorIs(t, service.ValidateResetToken(tok, testNow.Add(5*time.Minute)), domain.ErrResetTokenUsed)
}

// Expired проверяется раньше Used — если оба, побеждает Expired.
func TestValidateResetToken_ExpiredBeforeUsed(t *testing.T) {
	usedAt := testNow
	tok := makeResetToken(testNow, 30*time.Minute, &usedAt)
	assert.ErrorIs(t, service.ValidateResetToken(tok, testNow.Add(time.Hour)), domain.ErrResetTokenExpired)
}
