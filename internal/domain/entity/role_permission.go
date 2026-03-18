package entity

import (
	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

func NewRolePermission(roleID, permissionID uuid.UUID) (RolePermission, error) {
	if roleID == uuid.Nil || permissionID == uuid.Nil {
		return RolePermission{}, domain.ErrInvalidRolePermission
	}

	return RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}, nil
}
