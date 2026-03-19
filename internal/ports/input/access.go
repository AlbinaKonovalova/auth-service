package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// AssignRoleInput — входные данные для назначения роли пользователю.
type AssignRoleInput struct {
	UserID   uuid.UUID
	RoleCode string
}

// RevokeRoleInput — входные данные для отзыва роли у пользователя.
type RevokeRoleInput struct {
	UserID   uuid.UUID
	RoleCode string
}

// AssignPermissionInput — входные данные для назначения permission роли.
type AssignPermissionInput struct {
	RoleCode       string
	PermissionCode string
}

// RevokePermissionInput — входные данные для отзыва permission у роли.
type RevokePermissionInput struct {
	RoleCode       string
	PermissionCode string
}

// AccessUseCase — входной контракт для access admin API.
//
// GetUserRoles и GetRolePermissions возвращают доменный view —
// маппинг в transport-тип выполняется в controller.
type AccessUseCase interface {
	// GetUserRoles возвращает список ролей пользователя.
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]service.RoleView, error)

	// GetRolePermissions возвращает список permissions роли по её code.
	// Возвращает domain.ErrInvalidRoleCode если code невалиден.
	// Возвращает domain.ErrRoleNotFound если роль не существует.
	// Возвращает domain.ErrDataIntegrityViolation если обнаружена битая связь role ↔ permission.
	GetRolePermissions(ctx context.Context, roleCode string) ([]service.PermissionView, error)

	// AssignPermission назначает permission роли.
	// Идемпотентен: если permission уже назначен роли — возвращает success без ошибки.
	// Возвращает domain.ErrInvalidRoleCode если role code невалиден.
	// Возвращает domain.ErrInvalidPermissionCode если permission code невалиден.
	// Возвращает domain.ErrRoleNotFound если роль не существует.
	// Возвращает domain.ErrPermissionNotFound если permission не существует.
	AssignPermission(ctx context.Context, in AssignPermissionInput) error

	// RevokePermission отзывает permission у роли.
	// Возвращает domain.ErrInvalidRoleCode если role code невалиден.
	// Возвращает domain.ErrInvalidPermissionCode если permission code невалиден.
	// Возвращает domain.ErrRoleNotFound если роль не существует.
	// Возвращает domain.ErrPermissionNotFound если permission не существует.
	// Возвращает domain.ErrRolePermissionNotFound если permission не назначен роли.
	RevokePermission(ctx context.Context, in RevokePermissionInput) error

	// AssignRole назначает роль пользователю.
	// Идемпотентен: если роль уже назначена — возвращает success без ошибки.
	// Возвращает domain.ErrUserNotFound если пользователь не существует.
	// Возвращает domain.ErrRoleNotFound если роль с таким code не существует.
	// Возвращает domain.ErrInvalidRoleCode если code невалиден.
	AssignRole(ctx context.Context, in AssignRoleInput) error

	// RevokeRole отзывает роль у пользователя.
	// Возвращает domain.ErrInvalidRoleCode если code невалиден.
	// Возвращает domain.ErrUserNotFound если пользователь не существует.
	// Возвращает domain.ErrRoleNotFound если роль с таким code не существует.
	// Возвращает domain.ErrUserRoleNotFound если роль не назначена пользователю.
	// Возвращает domain.ErrCannotRevokeLastRole если это последняя роль пользователя.
	RevokeRole(ctx context.Context, in RevokeRoleInput) error
}
