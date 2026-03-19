package service

import (
	"sort"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

// PermissionView — доменное представление permission.
type PermissionView struct {
	ID          uuid.UUID
	Code        string
	Description string
}

// BuildPermissionListResult преобразует []entity.Permission в []PermissionView.
// Сортирует по code — стабильный порядок business result зафиксирован здесь,
// а не делегируется ORDER BY в repo.
// Используется в ListPermissions сценарии.
func BuildPermissionListResult(permissions []entity.Permission) []PermissionView {
	result := make([]PermissionView, len(permissions))
	for i, p := range permissions {
		result[i] = PermissionView{
			ID:          p.ID,
			Code:        p.Code,
			Description: p.Description,
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result
}

// BuildRolePermissionsResult формирует итоговый список permissions роли:
//   - дедублицирует по ID (на случай дублей в данных)
//   - проверяет, что для каждого permission_id из rolePerms существует permission в справочнике;
//     если связь битая (permission_id есть в role_permissions, но нет в справочнике) —
//     возвращает domain.ErrDataIntegrityViolation, а не бизнесовый ErrPermissionNotFound,
//     чтобы не смешивать нарушение целостности данных с обычным "not found" по запросу пользователя
//   - сортирует по code для стабильного порядка ответа
//
// Используется в GetRolePermissions сценарии.
func BuildRolePermissionsResult(rolePerms []entity.RolePermission, permissions []entity.Permission) ([]PermissionView, error) {
	if len(rolePerms) == 0 {
		return []PermissionView{}, nil
	}

	permByID := make(map[uuid.UUID]entity.Permission, len(permissions))
	for _, p := range permissions {
		permByID[p.ID] = p
	}

	seen := make(map[uuid.UUID]struct{}, len(rolePerms))
	result := make([]PermissionView, 0, len(rolePerms))

	for _, rp := range rolePerms {
		if _, dup := seen[rp.PermissionID]; dup {
			continue
		}
		seen[rp.PermissionID] = struct{}{}

		p, ok := permByID[rp.PermissionID]
		if !ok {
			return nil, domain.ErrDataIntegrityViolation
		}

		result = append(result, PermissionView{
			ID:          p.ID,
			Code:        p.Code,
			Description: p.Description,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result, nil
}
