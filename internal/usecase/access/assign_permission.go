package access

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *AccessService) AssignPermission(ctx context.Context, in input.AssignPermissionInput) error {
	roleCode, err := value.NewRoleCode(in.RoleCode)
	if err != nil {
		return err
	}

	permCode, err := value.NewPermissionCode(in.PermissionCode)
	if err != nil {
		return err
	}

	return s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		role, err := s.roleRepo.FindByCode(txCtx, roleCode.String())
		if err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}

		perm, err := s.permRepo.FindByCode(txCtx, permCode.String())
		if err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}

		rp, err := entity.NewRolePermission(role.ID, perm.ID)
		if err != nil {
			return fmt.Errorf("assign permission: build role permission: %w", err)
		}

		if err := s.rolePermRepo.Assign(txCtx, rp); err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}

		return nil
	})
}
