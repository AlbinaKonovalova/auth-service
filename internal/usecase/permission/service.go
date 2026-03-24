package permission

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type PermissionService struct {
	permRepo     output.PermissionRepository
	rolePermRepo output.RolePermissionRepository
	uuidGen      output.UUIDGenerator
	tx           output.TxManager
}

func NewPermissionService(
	permRepo output.PermissionRepository,
	rolePermRepo output.RolePermissionRepository,
	uuidGen output.UUIDGenerator,
	tx output.TxManager,
) *PermissionService {
	return &PermissionService{
		permRepo:     permRepo,
		rolePermRepo: rolePermRepo,
		uuidGen:      uuidGen,
		tx:           tx,
	}
}
