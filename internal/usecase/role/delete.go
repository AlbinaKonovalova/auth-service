package role

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

// DeleteRole удаляет роль по code.
// Сценарий выполняется в одной транзакции (read-check-delete):
//  1. доменная валидация role code через value.NewRoleCode
//  2. внутри tx — поиск роли с блокировкой строки roles (FOR UPDATE)
//     это сериализует delete по самой роли; дополнительно в обычной схеме с FK
//     конкурентные операции, создающие новые связи на эту роль, не должны
//     проходить мимо текущей транзакции незаметно
//  3. внутри tx — проверка наличия user_roles связей через ExistsByRoleID
//  4. внутри tx — проверка наличия role_permissions связей через ExistsByRoleID
//  5. доменное правило "нельзя удалить роль со связями" делегируется в domain/service.ValidateRoleDeletion
//  6. внутри tx — удаление роли
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
