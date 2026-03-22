package permission

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

// DeletePermission удаляет permission по code.
// Сценарий выполняется в одной транзакции (read-check-delete):
//  1. доменная валидация permission code через value.NewPermissionCode
//  2. внутри tx — поиск permission с блокировкой строки (FOR UPDATE)
//     сериализует конкурентные delete и гонку с assign permission к роли
//  3. внутри tx — проверка наличия role_permissions связей через ExistsByPermissionID
//  4. доменное правило "нельзя удалить используемый permission" — domain/service.ValidatePermissionDeletion
//  5. внутри tx — удаление permission
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
