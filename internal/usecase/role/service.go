package role

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

// RoleService реализует input.RoleUseCase.
// Оркестрирует сценарии управления справочником ролей.
type RoleService struct {
	roleRepo output.RoleRepository
}

func NewRoleService(roleRepo output.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}
