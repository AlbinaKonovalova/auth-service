package access

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// RevokeRole отзывает роль у пользователя.
// Весь сценарий выполняется в одной транзакции:
//  1. доменная валидация role code
//  2. проверка существования пользователя
//  3. получение роли по code
//  4. чтение текущих назначений с блокировкой (FOR UPDATE) — фиксируем состояние в рамках tx
//  5. domain/service.CanRevokeUserRole: проверяет наличие назначения и правило "нельзя снять последнюю роль"
//  6. удаление связи user ↔ role; repo возвращает ErrUserRoleNotFound при 0 affected rows (защита от гонки)
//
// Блокировка в шаге 4 гарантирует, что параллельный RevokeRole по тому же пользователю
// будет ждать завершения текущей транзакции и увидит уже изменённое состояние.
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
