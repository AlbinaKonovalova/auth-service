package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestNewEmail_Valid(t *testing.T) {
	cases := []struct {
		raw      string
		expected string
	}{
		{"admin@example.com", "admin@example.com"},
		{"Admin@Example.COM", "admin@example.com"}, // нормализация к lowercase
		{"  user@domain.io  ", "user@domain.io"},   // trim пробелов
		{"a@b.co", "a@b.co"},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			e, err := value.NewEmail(tc.raw)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, e.String())
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty string", ""},
		{"only spaces", "   "},
		{"no at-sign", "invalidemail"},
		{"no domain", "user@"},
		{"no local", "@domain.com"},
		{"multiple at-signs", "a@@b.com"},
		{"domain starts with dot", "user@.example.com"},
		{"domain ends with dot", "user@example."},
		{"domain without dot", "user@localhost"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := value.NewEmail(tc.raw)
			assert.ErrorIs(t, err, domain.ErrInvalidEmail)
		})
	}
}
