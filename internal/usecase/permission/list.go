package permission

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func (s *PermissionService) ListPermissions(ctx context.Context) ([]domainservice.PermissionView, error) {
	permissions, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	return domainservice.BuildPermissionListResult(permissions), nil
}
