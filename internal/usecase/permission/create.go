package permission

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// CreatePermission создаёт новый permission.
// Сценарий выполняется в одной транзакции (read-check-write):
//  1. доменная валидация permission code через value.NewPermissionCode
//  2. доменный конструктор entity.NewPermission — валидирует description (обязателен), нормализует
//  3. внутри tx — проверка уникальности code; если code занят — domain.ErrDuplicatePermissionCode
//  4. внутри tx — вставка нового permission; при race unique violation в DB также domain.ErrDuplicatePermissionCode
//
// Итоговое представление формируется через domain/service.PermissionViewFromEntity —
// единственный источник истины для маппинга entity.Permission в доменный result.
func (s *PermissionService) CreatePermission(ctx context.Context, in input.CreatePermissionInput) (domainservice.PermissionView, error) {
	permCode, err := value.NewPermissionCode(in.Code)
	if err != nil {
		return domainservice.PermissionView{}, err
	}

	id := s.uuidGen.New()

	perm, err := entity.NewPermission(entity.NewPermissionParams{
		ID:          id,
		Code:        permCode.String(),
		Description: in.Description,
	})
	if err != nil {
		return domainservice.PermissionView{}, err
	}

	if err := s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := s.permRepo.ExistsByCode(txCtx, perm.Code)
		if err != nil {
			return fmt.Errorf("check permission code uniqueness: %w", err)
		}
		if exists {
			return domain.ErrDuplicatePermissionCode
		}

		if err := s.permRepo.Create(txCtx, perm); err != nil {
			return fmt.Errorf("create permission: %w", err)
		}

		return nil
	}); err != nil {
		return domainservice.PermissionView{}, err
	}

	return domainservice.PermissionViewFromEntity(perm), nil
}
