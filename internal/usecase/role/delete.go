package role

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func (s *RoleService) DeleteRole(ctx context.Context, roleCode string) error {
	code, err := value.NewRoleCode(roleCode)
	if err != nil {
		return err
	}

	return s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		role, err := s.roleRepo.FindByCodeForUpdate(txCtx, code.String())
		if err != nil {
			return fmt.Errorf("delete role: find role: %w", err)
		}

		hasUsers, err := s.userRoleRepo.ExistsByRoleID(txCtx, role.ID)
		if err != nil {
			return fmt.Errorf("delete role: check user assignments: %w", err)
		}

		hasPerms, err := s.rolePermRepo.ExistsByRoleID(txCtx, role.ID)
		if err != nil {
			return fmt.Errorf("delete role: check permission assignments: %w", err)
		}

		if err := domainservice.ValidateRoleDeletion(hasUsers, hasPerms); err != nil {
			return err
		}

		if err := s.roleRepo.Delete(txCtx, role.ID); err != nil {
			return fmt.Errorf("delete role: %w", err)
		}

		return nil
	})
}
