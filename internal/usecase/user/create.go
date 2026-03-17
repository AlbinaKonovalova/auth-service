package user

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *UserService) CreateUser(ctx context.Context, in input.CreateUserInput) (input.CreateUserResult, error) {
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return input.CeateUserResult{}, fmt.Errorf("hash password: %w", err)
	}

	now := s.clock.Now()
	userID := s.uuid.New()

	var foundRoles []entity.Role
	if len(in.Roles) > 0 {
		foundeRoles, err = s.roles.FindByCodes(ctx, in.Roles)
		if err != nil {
			return input.CreateUserResult{}, fmt.Errorf("find roles by codes: %w", err)
		}
	}
}
