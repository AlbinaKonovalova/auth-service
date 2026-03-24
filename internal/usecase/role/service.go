package role

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type RoleService struct {
	roleRepo     output.RoleRepository
	userRoleRepo output.UserRoleRepository
	rolePermRepo output.RolePermissionRepository
	uuidGen      output.UUIDGenerator
	tx           output.TxManager
}

func NewRoleService(
	roleRepo output.RoleRepository,
	userRoleRepo output.UserRoleRepository,
	rolePermRepo output.RolePermissionRepository,
	uuidGen output.UUIDGenerator,
	tx output.TxManager,
) *RoleService {
	return &RoleService{
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		rolePermRepo: rolePermRepo,
		uuidGen:      uuidGen,
		tx:           tx,
	}
}
