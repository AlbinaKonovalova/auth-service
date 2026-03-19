package entity

import (
	"strings"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
)

type Role struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

// NewRoleParams — входные данные для создания новой роли.
type NewRoleParams struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
}

// NewRole создаёт валидную роль.
// Валидирует:
//   - code: обязателен, trimmed, lowercased (делегируется value.RoleCode снаружи)
//   - name: обязателен, непустой после trim
//   - description: необязателен
//
// Конструктор принимает уже нормализованный code (после value.NewRoleCode),
// чтобы не дублировать валидацию.
func NewRole(p NewRoleParams) (Role, error) {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return Role{}, domain.ErrRoleNameEmpty
	}

	return Role{
		ID:          p.ID,
		Code:        p.Code,
		Name:        name,
		Description: strings.TrimSpace(p.Description),
	}, nil
}
