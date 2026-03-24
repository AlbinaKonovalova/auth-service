package permission

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func (s *PermissionService) DeletePermission(ctx context.Context, permissionCode string) error {
	code, err := value.NewPermissionCode(permissionCode)
	if err != nil {
		return err
	}

	return s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		perm, err := s.permRepo.FindByCodeForUpdate(txCtx, code.String())
		if err != nil {
			return fmt.Errorf("delete permission: find permission: %w", err)
		}

		hasRoles, err := s.rolePermRepo.ExistsByPermissionID(txCtx, perm.ID)
		if err != nil {
			return fmt.Errorf("delete permission: check role assignments: %w", err)
		}

		if err := domainservice.ValidatePermissionDeletion(hasRoles); err != nil {
			return err
		}

		if err := s.permRepo.Delete(txCtx, perm.ID); err != nil {
			return fmt.Errorf("delete permission: %w", err)
		}

		return nil
	})
}
