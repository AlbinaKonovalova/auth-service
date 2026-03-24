package entity

import (
	"strings"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type Permission struct {
	ID          uuid.UUID
	Code        string
	Description string
}

// NewPermissionParams — параметры для создания новой permission.
type NewPermissionParams struct {
	ID          uuid.UUID
	Code        string
	Description string
}

func NewPermission(p NewPermissionParams) (Permission, error) {
	description := strings.TrimSpace(p.Description)
	if description == "" {
		return Permission{}, domain.ErrPermissionDescriptionEmpty
	}

	return Permission{
		ID:          p.ID,
		Code:        p.Code,
		Description: description,
	}, nil
}
