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

var testNow = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func TestMapLoginLookupError_UserNotFound(t *testing.T) {
	err := service.MapLoginLookupError(domain.ErrUserNotFound)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestMapLoginLookupError_OtherError_NotMasked(t *testing.T) {
	original := domain.ErrDataIntegrityViolation
	err := service.MapLoginLookupError(original)
	assert.ErrorIs(t, err, original)
}

func TestCanLogout_SessionNotFound_SilentSuccess(t *testing.T) {
	canRevoke, err := service.CanLogout(domain.ErrRefreshTokenNotFound, nil)
	assert.NoError(t, err)
	assert.False(t, canRevoke)
}

func TestCanLogout_InfraError_Propagated(t *testing.T) {
	infraErr := domain.ErrDataIntegrityViolation
	canRevoke, err := service.CanLogout(infraErr, nil)
	assert.ErrorIs(t, err, infraErr)
	assert.False(t, canRevoke)
}

func TestCanLogout_NilSession(t *testing.T) {
	canRevoke, err := service.CanLogout(nil, nil)
	assert.NoError(t, err)
	assert.False(t, canRevoke)
}

func TestCanLogout_RevokedSession(t *testing.T) {
	revokedAt := testNow
	s := &entity.RefreshSession{RevokedAt: &revokedAt}
	canRevoke, err := service.CanLogout(nil, s)
	assert.NoError(t, err)
	assert.False(t, canRevoke)
}

func TestCanLogout_ActiveSession(t *testing.T) {
	s := &entity.RefreshSession{ID: uuid.New()}
	canRevoke, err := service.CanLogout(nil, s)
	assert.NoError(t, err)
	assert.True(t, canRevoke)
}
