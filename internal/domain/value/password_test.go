package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func TestNewPassword_Valid(t *testing.T) {
	cases := []string{
		"12345678", // ровно 8 символов — граница
		"ValidPassword1!",
		"a very long password that is definitely valid and secure",
	}

	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			p, err := value.NewPassword(raw)
			require.NoError(t, err)
			assert.Equal(t, raw, p.String())
		})
	}
}

func TestNewPassword_Invalid(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"empty", ""},
		{"7 chars — just below minimum", "1234567"},
		{"single char", "x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := value.NewPassword(tc.raw)
			assert.ErrorIs(t, err, domain.ErrInvalidPassword)
		})
	}
}

func TestNewPassword_MinLength(t *testing.T) {
	// граничный случай: ровно MinPasswordLength символов — должен пройти
	exact := make([]byte, value.MinPasswordLength)
	for i := range exact {
		exact[i] = 'a'
	}
	_, err := value.NewPassword(string(exact))
	assert.NoError(t, err)

	// на один меньше — должен упасть
	short := make([]byte, value.MinPasswordLength-1)
	for i := range short {
		short[i] = 'a'
	}
	_, err = value.NewPassword(string(short))
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}
