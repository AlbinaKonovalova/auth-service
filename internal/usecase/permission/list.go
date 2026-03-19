package permission

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// ListPermissions возвращает полный список permissions в стабильном порядке (по code).
// Сортировка зафиксирована в domain/service.BuildPermissionListResult — это источник истины для порядка результата.
func (s *PermissionService) ListPermissions(ctx context.Context) ([]domainservice.PermissionView, error) {
	permissions, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	return domainservice.BuildPermissionListResult(permissions), nil
}
