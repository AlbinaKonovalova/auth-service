package access

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *AccessService) RevokeRole(ctx context.Context, in input.RevokeRoleInput) error {
	roleCode, err := value.NewRoleCode(in.RoleCode)
	if err != nil {
		return err
	}

	return s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := s.userRepo.FindByID(txCtx, in.UserID)
		if err != nil {
			return fmt.Errorf("revoke role: %w", err)
		}

		role, err := s.roleRepo.FindByCode(txCtx, roleCode.String())
		if err != nil {
			return fmt.Errorf("revoke role: %w", err)
		}

		allRoles, err := s.userRoleRepo.FindByUserIDForUpdate(txCtx, in.UserID)
		if err != nil {
			return fmt.Errorf("revoke role: find user roles: %w", err)
		}

		if err := domainservice.CanRevokeUserRole(allRoles, role.ID); err != nil {
			return err
		}

		if err := s.userRoleRepo.Revoke(txCtx, in.UserID, role.ID); err != nil {
			return fmt.Errorf("revoke role: %w", err)
		}

		return nil
	})
}
