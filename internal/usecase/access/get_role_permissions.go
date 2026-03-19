package access

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

// GetRolePermissions возвращает список permissions роли.
// Порядок:
//  1. доменная валидация role code
//  2. получение роли по code — возвращает ErrRoleNotFound если роль не существует
//  3. чтение связей role ↔ permission
//  4. загрузка permission справочника по найденным IDs
//  5. сборка итогового доменного результата — domain/service.BuildRolePermissionsResult:
//     проверяет битые связи, дедуплицирует, сортирует по code
func (s *AccessService) GetRolePermissions(ctx context.Context, roleCode string) ([]domainservice.PermissionView, error) {
	rc, err := value.NewRoleCode(roleCode)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByCode(ctx, rc.String())
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
	}

	rolePerms, err := s.rolePermRepo.FindByRoleID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: find role permissions: %w", err)
	}

	if len(rolePerms) == 0 {
		return []domainservice.PermissionView{}, nil
	}

	seen := make(map[uuid.UUID]struct{}, len(rolePerms))
	permIDs := make([]uuid.UUID, 0, len(rolePerms))
	for _, rp := range rolePerms {
		if _, ok := seen[rp.PermissionID]; !ok {
			seen[rp.PermissionID] = struct{}{}
			permIDs = append(permIDs, rp.PermissionID)
		}
	}

	permissions, err := s.permRepo.FindByIDs(ctx, permIDs)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: find permissions: %w", err)
	}

	return domainservice.BuildRolePermissionsResult(rolePerms, permissions)
}
