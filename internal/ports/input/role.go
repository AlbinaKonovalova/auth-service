package input

import (
	"context"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// RoleUseCase — входной контракт для roles admin API.
type RoleUseCase interface {
	// ListRoles возвращает полный список ролей.
	ListRoles(ctx context.Context) ([]domainservice.RoleView, error)
}
