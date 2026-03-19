package input

import (
	"context"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// CreateRoleInput — входные данные для создания новой роли.
type CreateRoleInput struct {
	Code        string
	Name        string
	Description string
}

// RoleUseCase — входной контракт для roles admin API.
type RoleUseCase interface {
	// ListRoles возвращает полный список ролей.
	ListRoles(ctx context.Context) ([]domainservice.RoleView, error)

	// CreateRole создаёт новую роль.
	// Возвращает domain.ErrInvalidRoleCode если code невалиден.
	// Возвращает domain.ErrRoleNameEmpty если name пустой.
	// Возвращает domain.ErrDuplicateRoleCode если роль с таким code уже существует.
	CreateRole(ctx context.Context, in CreateRoleInput) (domainservice.RoleView, error)

	// DeleteRole удаляет роль по code.
	// Возвращает domain.ErrInvalidRoleCode если code невалиден.
	// Возвращает domain.ErrRoleNotFound если роль не существует.
	// Возвращает domain.ErrRoleInUse если роль назначена пользователям или имеет permissions.
	DeleteRole(ctx context.Context, roleCode string) error
}
