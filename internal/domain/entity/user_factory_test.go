package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestBuildNewUserAggregate_Valid(t *testing.T) {
	adminRole := entity.Role{ID: uuid.New(), Code: "admin", Name: "Admin"}

	agg, err := entity.BuildNewUserAggregate(
		testUserID,
		"user@example.com",
		"ValidPass1!",
		"hashed_password",
		[]string{"admin"},
		[]entity.Role{adminRole},
		testNow,
	)

	require.NoError(t, err)
	assert.Equal(t, testUserID, agg.User.ID)
	assert.True(t, agg.User.IsActive)
	assert.Len(t, agg.UserRoles, 1)
	assert.Equal(t, adminRole.ID, agg.UserRoles[0].RoleID)
}

func TestBuildNewUserAggregate_InvalidEmail(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "not-an-email", "ValidPass1!", "hash",
		[]string{"admin"}, []entity.Role{{Code: "admin"}}, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)
}

func TestBuildNewUserAggregate_InvalidPassword(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "short", "hash",
		[]string{"admin"}, []entity.Role{{Code: "admin"}}, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}

func TestBuildNewUserAggregate_NoRoles(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{}, nil, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrUserMustHaveRole)
}

func TestBuildNewUserAggregate_DuplicateRoleCodes(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{"admin", "ADMIN"}, nil, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrDuplicateRoleCode)
}

func TestBuildNewUserAggregate_RoleNotFound(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{"admin"}, []entity.Role{{Code: "manager"}}, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrRoleNotFound)
}

func TestBuildNewUserAggregate_MultipleRoles(t *testing.T) {
	role1 := entity.Role{ID: uuid.New(), Code: "admin", Name: "Admin"}
	role2 := entity.Role{ID: uuid.New(), Code: "manager", Name: "Manager"}

	agg, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{"admin", "manager"},
		[]entity.Role{role1, role2},
		testNow,
	)
	require.NoError(t, err)
	assert.Len(t, agg.UserRoles, 2)
}

func TestBuildNewUserAggregate_NormalizeEmailToLower(t *testing.T) {
	role := entity.Role{ID: uuid.New(), Code: "admin", Name: "Admin"}
	agg, err := entity.BuildNewUserAggregate(
		testUserID, "User@Example.COM", "ValidPass1!", "hash",
		[]string{"admin"}, []entity.Role{role}, testNow,
	)
	require.NoError(t, err)
	assert.Equal(t, "user@example.com", agg.User.Email)
}

func TestBuildNewUserAggregate_InvalidRoleCode(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{""}, nil, testNow,
	)
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestBuildNewUserAggregate_RoleCodeNormalization(t *testing.T) {
	role := entity.Role{ID: uuid.New(), Code: "admin", Name: "Admin"}

	agg, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{"ADMIN"}, []entity.Role{role}, testNow,
	)
	require.NoError(t, err)
	assert.Len(t, agg.UserRoles, 1)
}

func TestBuildNewUserAggregate_RequestedRoleNotInFoundList(t *testing.T) {
	_, err := entity.BuildNewUserAggregate(
		testUserID, "user@example.com", "ValidPass1!", "hash",
		[]string{"superadmin"},
		[]entity.Role{{ID: uuid.New(), Code: "admin", Name: "Admin"}},
		testNow,
	)
	assert.ErrorIs(t, err, domain.ErrRoleNotFound)
}

func newValidRole(code string) entity.Role {
	return entity.Role{ID: uuid.New(), Code: code, Name: code}
}

var _ = value.Email{} // suppress unused import warning if needed
