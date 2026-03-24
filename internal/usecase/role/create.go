package role

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *RoleService) CreateRole(ctx context.Context, in input.CreateRoleInput) (domainservice.RoleView, error) {
	roleCode, err := value.NewRoleCode(in.Code)
	if err != nil {
		return domainservice.RoleView{}, err
	}

	id := s.uuidGen.New()

	role, err := entity.NewRole(entity.NewRoleParams{
		ID:          id,
		Code:        roleCode.String(),
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		return domainservice.RoleView{}, err
	}

	if err := s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := s.roleRepo.ExistsByCode(txCtx, role.Code)
		if err != nil {
			return fmt.Errorf("check role code uniqueness: %w", err)
		}
		if exists {
			return domain.ErrDuplicateRoleCode
		}

		if err := s.roleRepo.Create(txCtx, role); err != nil {
			return fmt.Errorf("create role: %w", err)
		}

		return nil
	}); err != nil {
		return domainservice.RoleView{}, err
	}

	return domainservice.RoleViewFromEntity(role), nil
}
