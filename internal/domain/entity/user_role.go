package entity

import (
	"time"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/google/uuid"
)

type UserRole struct {
	UserID    uuid.UUID
	RoleID    uuid.UUID
	CreatedAt time.Time
}

func NewUserRole(userID, roleID uuid.UUID, now time.Time) (UserRole, error) {
	if userID == uuid.Nil || roleID == uuid.Nil || now.IsZero() {
		return UserRole{}, domain.ErrInvalidUserRole
	}

	return UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: now,
	}, nil
}
