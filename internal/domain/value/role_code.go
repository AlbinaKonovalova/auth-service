package value

import (
	"strings"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type RoleCode sruct {
	value string
}

func NewRoleCode(raw string) (RoleCode, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return RooleCode{}, domain.ErrInvalidRoleCode}
	}

	return RoleCode{value: v}, nil
}

func (r RoleCode) String() string {
	return r.value
}