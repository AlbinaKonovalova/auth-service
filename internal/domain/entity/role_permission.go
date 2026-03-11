package entity

import "github.com/google/uuid"

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}
