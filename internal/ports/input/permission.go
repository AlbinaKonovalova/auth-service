package input

import (
	"context"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

type CreatePermissionInput struct {
	Code        string
	Description string
}

type PermissionUseCase interface {
	ListPermissions(ctx context.Context) ([]domainservice.PermissionView, error)

	CreatePermission(ctx context.Context, in CreatePermissionInput) (domainservice.PermissionView, error)

	DeletePermission(ctx context.Context, permissionCode string) error
}
