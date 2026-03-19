package access

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// AssignRole назначает роль пользователю.
// Весь сценарий выполняется в одной транзакции:
//  1. доменная валидация role code
//  2. проверка существования пользователя
//  3. получение роли по code
//  4. явная проверка: роль уже назначена → вернуть nil (идемпотентность на уровне сценария)
//  5. создание entity.UserRole
//  6. сохранение; ON CONFLICT DO NOTHING в repo — второй слой защиты от гонок
func (s *AccessService) AssignRole(ctx context.Context, in input.AssignRoleInput) error {
	roleCode, err := value.NewRoleCode(in.RoleCode)
	if err != nil {
		return err
	}

	return s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		_, err := s.userRepo.FindByID(txCtx, in.UserID)
		if err != nil {
			return fmt.Errorf("assign role: %w", err)
		}

		role, err := s.roleRepo.FindByCode(txCtx, roleCode.String())
		if err != nil {
			return fmt.Errorf("assign role: %w", err)
		}

		exists, err := s.userRoleRepo.Exists(txCtx, in.UserID, role.ID)
		if err != nil {
			return fmt.Errorf("assign role: check exists: %w", err)
		}
		if exists {
			return nil
		}

		userRole, err := entity.NewUserRole(in.UserID, role.ID, s.clock.Now())
		if err != nil {
			return fmt.Errorf("assign role: build user role: %w", err)
		}

		if err := s.userRoleRepo.Assign(txCtx, userRole); err != nil {
			return fmt.Errorf("assign role: %w", err)
		}

		return nil
	})
}
