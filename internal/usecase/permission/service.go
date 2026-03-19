package permission

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

// PermissionService реализует input.PermissionUseCase.
// Оркестрирует сценарии управления справочником permissions.
type PermissionService struct {
	permRepo output.PermissionRepository
}

func NewPermissionService(permRepo output.PermissionRepository) *PermissionService {
	return &PermissionService{permRepo: permRepo}
}
