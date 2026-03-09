package dto

import "github.com/google/uuid"

type UserAccess struct {
	ID          uuid.UUID
	Email       string
	IsActive    bool
	Roles       []string
	Permissions []string
}
