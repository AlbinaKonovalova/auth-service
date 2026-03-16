package value

import (
	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

const MinPasswordLength = 8

type Password struct {
	value string
}

func NewPassword(raw string) (Password, error) {
	if len(raw) < MinPasswordLength {
		return Password{}, domain.ErrInvalidPassword
	}
	return Password{value: raw}, nil
}

func (p Password) String() string {
	return p.value
}
