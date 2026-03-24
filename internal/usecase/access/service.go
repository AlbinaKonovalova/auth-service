package access

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type AccessService struct {
	userRepo     output.UserRepository
	userRoleRepo output.UserRoleRepository
	roleRepo     output.RoleRepository
	rolePermRepo output.RolePermissionRepository
	permRepo     output.PermissionRepository
	clock        output.Clock
	tx           output.TxManager
}

func NewAccessService(
	userRepo output.UserRepository,
	userRoleRepo output.UserRoleRepository,
	roleRepo output.RoleRepository,
	rolePermRepo output.RolePermissionRepository,
	permRepo output.PermissionRepository,
	clock output.Clock,
	tx output.TxManager,
) *AccessService {
	return &AccessService{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		rolePermRepo: rolePermRepo,
		permRepo:     permRepo,
		clock:        clock,
		tx:           tx,
	}
}
