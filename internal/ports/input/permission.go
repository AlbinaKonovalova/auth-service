package input

import (
	"context"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// PermissionUseCase — входной контракт для permissions admin API.
type PermissionUseCase interface {
	// ListPermissions возвращает полный список permissions.
	// Возвращает доменный []service.PermissionView — маппинг в transport выполняется в controller.
	ListPermissions(ctx context.Context) ([]domainservice.PermissionView, error)
}
