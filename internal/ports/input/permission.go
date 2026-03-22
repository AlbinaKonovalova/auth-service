package input

import (
	"context"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// CreatePermissionInput — входные данные для создания нового permission.
type CreatePermissionInput struct {
	Code        string
	Description string
}

// PermissionUseCase — входной контракт для permissions admin API.
type PermissionUseCase interface {
	// ListPermissions возвращает полный список permissions.
	// Возвращает доменный []service.PermissionView — маппинг в transport выполняется в controller.
	ListPermissions(ctx context.Context) ([]domainservice.PermissionView, error)

	// CreatePermission создаёт новый permission.
	// Возвращает domain.ErrInvalidPermissionCode если code невалиден.
	// Возвращает domain.ErrPermissionDescriptionEmpty если description пустой после trim.
	// Возвращает domain.ErrDuplicatePermissionCode если permission с таким code уже существует.
	CreatePermission(ctx context.Context, in CreatePermissionInput) (domainservice.PermissionView, error)

	// DeletePermission удаляет permission по code.
	// Возвращает domain.ErrInvalidPermissionCode если code невалиден.
	// Возвращает domain.ErrPermissionNotFound если permission не существует.
	// Возвращает domain.ErrPermissionInUse если permission назначен хотя бы одной роли.
	DeletePermission(ctx context.Context, permissionCode string) error
}
