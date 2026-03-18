package input

import (
	"context"

	"github.com/google/uuid"
)

// RoleResult — представление роли в ответах access API.
type RoleResult struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

// PermissionResult — представление permission в ответах access API.
type PermissionResult struct {
	ID          uuid.UUID
	Code        string
	Description string
}

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

// AccessUseCase — входной контракт для access admin API.
//
// Вычисление effective roles / permissions пользователя (dedup + stable order)
// является доменным результатом и реализовано в domain/service.ResolveEffectiveAccess.
// Usecase загружает данные через output ports и делегирует сборку туда.
type AccessUseCase interface {
	// GetUserRoles возвращает список ролей пользователя.
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleResult, error)

	// AssignRole назначает роль пользователю.
	// Идемпотентен: если роль уже назначена — возвращает success.
	AssignRole(ctx context.Context, in AssignRoleInput) error

	// RevokeRole отзывает роль у пользователя.
	// Возвращает ошибку если:
	//   - роль не назначена пользователю (ErrUserRoleNotFound)
	//   - это последняя роль пользователя (ErrCannotRevokeLastRole)
	RevokeRole(ctx context.Context, in RevokeRoleInput) error

	// ListRoles возвращает полный список ролей.
	ListRoles(ctx context.Context) ([]RoleResult, error)

	// ListPermissions возвращает полный список permissions.
	ListPermissions(ctx context.Context) ([]PermissionResult, error)
}
