package user

import (
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type UserService struct {
	users     output.UserRepository
	roles     output.RoleRepository
	userRoles output.UserRoleRepository
	sessions  output.RefreshSessionRepository
	hasher    output.PasswordHasher
	tx        output.TxManager
	clock     output.Clock
	uuid      output.UUIDGenerator
}

func NewUserService(
	users output.UserRepository,
	roles output.RoleRepository,
	userRoles output.UserRoleRepository,
	sessions output.RefreshSessionRepository,
	hasher output.PasswordHasher,
	tx output.TxManager,
	clock output.Clock,
	uuid output.UUIDGenerator,
) *UserService {
	return &UserService{
		users:     users,
		roles:     roles,
		userRoles: userRoles,
		sessions:  sessions,
		hasher:    hasher,
		tx:        tx,
		clock:     clock,
		uuid:      uuid,
	}
}
