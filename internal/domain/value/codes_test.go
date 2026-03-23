package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestNewRoleCode_Valid(t *testing.T) {
	cases := []struct {
		raw      string
		expected string
	}{
		{"admin", "admin"},
		{"ADMIN", "admin"},         // нормализация к lowercase
		{"  manager  ", "manager"}, // trim
		{"content_maker", "content_maker"},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			rc, err := value.NewRoleCode(tc.raw)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, rc.String())
		})
	}
}

func TestNewRoleCode_Invalid(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty string", ""},
		{"only spaces", "   "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := value.NewRoleCode(tc.raw)
			assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
		})
	}
}

func TestNewPermissionCode_Valid(t *testing.T) {
	cases := []struct {
		raw      string
		expected string
	}{
		{"users.read", "users.read"},
		{"USERS.READ", "users.read"},
		{"  stats.export  ", "stats.export"},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			pc, err := value.NewPermissionCode(tc.raw)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, pc.String())
		})
	}
}

func TestNewPermissionCode_Invalid(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty string", ""},
		{"only spaces", "   "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := value.NewPermissionCode(tc.raw)
			assert.ErrorIs(t, err, domain.ErrInvalidPermissionCode)
		})
	}
}
