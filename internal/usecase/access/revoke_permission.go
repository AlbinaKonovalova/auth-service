package access

import (
	"context"
	"fmt"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// RevokePermission отзывает permission у роли.
// Весь сценарий выполняется в одной транзакции:
//  1. доменная валидация role code и permission code
//  2. получение роли по code — возвращает ErrRoleNotFound если не существует
//  3. получение permission по code — возвращает ErrPermissionNotFound если не существует
//  4. проверка существования связи role ↔ permission — возвращает ErrRolePermissionNotFound если связи нет
//  5. удаление связи
func (s *AccessService) RevokePermission(ctx context.Context, in input.RevokePermissionInput) error {
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
			return fmt.Errorf("revoke permission: %w", err)
		}

		perm, err := s.permRepo.FindByCode(txCtx, permCode.String())
		if err != nil {
			return fmt.Errorf("revoke permission: %w", err)
		}

		exists, err := s.rolePermRepo.Exists(txCtx, role.ID, perm.ID)
		if err != nil {
			return fmt.Errorf("revoke permission: check exists: %w", err)
		}
		if !exists {
			return domain.ErrRolePermissionNotFound
		}

		if err := s.rolePermRepo.Revoke(txCtx, role.ID, perm.ID); err != nil {
			return fmt.Errorf("revoke permission: %w", err)
		}

		return nil
	})
}
