package service

import (
	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CanRevokeUserRole проверяет доменные правила отзыва роли у пользователя:
//  1. роль должна быть назначена пользователю → иначе ErrUserRoleNotFound
//  2. нельзя снять последнюю роль у пользователя → иначе ErrCannotRevokeLastRole
//
// Принимает все текущие назначения пользователя и roleID который хотим снять.
func CanRevokeUserRole(userRoles []entity.UserRole, revokeRoleID uuid.UUID) error {
	assigned := false
	remaining := 0

	for _, ur := range userRoles {
		if ur.RoleID == revokeRoleID {
			assigned = true
		} else {
			remaining++
		}
	}

	if !assigned {
		return domain.ErrUserRoleNotFound
	}
	if remaining == 0 {
		return domain.ErrCannotRevokeLastRole
	}
	return nil
}
