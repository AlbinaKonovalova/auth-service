package entity

import (
	"strings"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type Role struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

type NewRoleParams struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

func NewRole(p NewRoleParams) (Role, error) {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return Role{}, domain.ErrRoleNameEmpty
	}

	return Role{
		ID:          p.ID,
		Code:        p.Code,
		Name:        name,
		Description: strings.TrimSpace(p.Description),
	}, nil
}
