package dto

import "github.com/google/uuid"

type CurrentUser struct {
	ID          uuid.UUID
	Email       string
	IsActive    bool
	Roles       []string
	Permissions []string
}
