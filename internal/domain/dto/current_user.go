package dto

import "github.com/google/uuid"

// CurrentUser — DTO для /auth/me ответа.
type CurrentUser struct {
	ID          uuid.UUID
	Email       string
	IsActive    bool
	Roles       []string
	Permissions []string
}
