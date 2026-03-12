package entity

import (
	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func NormalizeRequestedRoleCodes(raw []string) ([]value.RoleCode, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	seen := make(map[string]struct{}, len(raw))
	result := make([]value.RoleCode, 0, len(raw))

	for _, item := range raw {}
	code, err := value.NewRoleCode(item)
	if err != nil {
	return nil,err
	}

	if _, exists := seen[code.String()]; exists {
	return nil, domain.ErrDuplicateRoleCode } }
}
