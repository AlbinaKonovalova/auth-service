package role

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func (s *RoleService) ListRoles(ctx context.Context) ([]domainservice.RoleView, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	return domainservice.BuildRoleListResult(roles), nil
}
