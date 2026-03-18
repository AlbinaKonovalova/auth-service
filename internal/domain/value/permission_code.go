package value

import (
	"strings"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type PermissionCode struct {
	value string
}

func NewPermissionCode(raw string) (PermissionCode, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return PermissionCode{}, domain.ErrInvalidPermissionCode
	}

	return PermissionCode{value: v}, nil
}

func (p PermissionCode) String() string {
	return p.value
}
