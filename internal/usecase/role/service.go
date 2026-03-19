package role

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

// RoleService реализует input.RoleUseCase.
// Оркестрирует сценарии управления справочником ролей.
type RoleService struct {
	roleRepo output.RoleRepository
	uuidGen  output.UUIDGenerator
	tx       output.TxManager
}

func NewRoleService(roleRepo output.RoleRepository, uuidGen output.UUIDGenerator, tx output.TxManager) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		uuidGen:  uuidGen,
		tx:       tx,
	}
}
