package access

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// AssignPermission назначает permission роли.
// Весь сценарий выполняется в одной транзакции:
//  1. доменная валидация role code и permission code
//  2. получение роли по code — возвращает ErrRoleNotFound если не существует
//  3. получение permission по code — возвращает ErrPermissionNotFound если не существует
//  4. создание entity.RolePermission
//  5. сохранение; ON CONFLICT DO NOTHING в repo обеспечивает идемпотентность при гонке —
//     usecase задаёт семантику сценария, база атомарно исполняет уже принятое решение
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
