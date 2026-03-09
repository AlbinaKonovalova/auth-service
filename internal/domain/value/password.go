package value

import "errors"

const passwordMinLen = 8

type Password struct {
	value string
}

func NewPassword(raw string) (Password, error) {
	if len(raw) < passwordMinLen {
		return Password{}, errors.New("password must be at least 8 characters")
	}
	return Password{value: raw}, nil
}

func (p Password) String() string {
	return p.value
}
