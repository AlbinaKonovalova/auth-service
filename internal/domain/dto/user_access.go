package dto

import "github.com/google/uuid"

// UserAccess — DTO пользователя вместе с ролями и permissions.
type UserAccess struct {
	ID          uuid.UUID
	Email       string
	IsActive    bool
	Roles       []string
	Permissions []string
}
