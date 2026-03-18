package value

import (
	"strings"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return Email{}, domain.ErrInvalidEmail
	}

	parts := strings.Split(v, "@")
	if len(parts) != 2 {
		return Email{}, domain.ErrInvalidEmail
	}

	local := parts[0]
	domainPart := parts[1]

	if local == "" || domainPart == "" {
		return Email{}, domain.ErrInvalidEmail
	}

	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return Email{}, domain.ErrInvalidEmail
	}

	if !strings.Contains(domainPart, ".") {
		return Email{}, domain.ErrInvalidEmail
	}

	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}
