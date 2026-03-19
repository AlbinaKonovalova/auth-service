package service

import (
	"sort"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// RoleView — доменное представление роли.
type RoleView struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

// BuildRoleListResult преобразует []entity.Role в []RoleView.
// Сортирует по code — стабильный порядок business result зафиксирован здесь,
// а не делегируется ORDER BY в repo.
// Используется в ListRoles сценарии.
func BuildRoleListResult(roles []entity.Role) []RoleView {
	result := make([]RoleView, len(roles))
	for i, r := range roles {
		result[i] = RoleView{
			ID:          r.ID,
			Code:        r.Code,
			Name:        r.Name,
			Description: r.Description,
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})
	return result
}

// BuildUserRolesResult формирует итоговый список ролей пользователя:
//   - дедублицирует по ID (на случай дублей в данных)
//   - проверяет, что для каждого role_id из userRoles существует роль в справочнике;
//     если связь битая (role_id есть в user_roles, но нет в справочнике) —
//     возвращает domain.ErrDataIntegrityViolation, а не бизнесовый ErrRoleNotFound,
//     чтобы не смешивать нарушение целостности данных с обычным "not found" по запросу пользователя
//   - сортирует по code для стабильного порядка ответа
func BuildUserRolesResult(userRoles []entity.UserRole, roles []entity.Role) ([]RoleView, error) {
	if len(userRoles) == 0 {
		return []RoleView{}, nil
	}

	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, r := range roles {
		roleByID[r.ID] = r
	}

	seen := make(map[uuid.UUID]struct{}, len(userRoles))
	result := make([]RoleView, 0, len(userRoles))

	for _, ur := range userRoles {
		if _, dup := seen[ur.RoleID]; dup {
			continue
		}
		seen[ur.RoleID] = struct{}{}

		r, ok := roleByID[ur.RoleID]
		if !ok {
			return nil, domain.ErrDataIntegrityViolation
		}

		result = append(result, RoleView{
			ID:          r.ID,
			Code:        r.Code,
			Name:        r.Name,
			Description: r.Description,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result, nil
}
