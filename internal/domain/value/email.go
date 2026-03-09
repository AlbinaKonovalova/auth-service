package value

import (
	"errors"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return Email{}, errors.New("email must not be empty")
	}

	parts := strings.Split(v, "@")
	if len(parts) != 2 {
		return Email{}, errors.New("invalid email format")
	}

	local := parts[0]
	domain := parts[1]

	if local == "" || domain == "" {
		return Email{}, errors.New("invalid email format")
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return Email{}, errors.New("invalid email format")
	}

	if !strings.Contains(domain, ".") {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}
