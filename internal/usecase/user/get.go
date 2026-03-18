package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (input.GetUserResult, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return input.GetUserResult{}, fmt.Errorf("find user by id: %w", err)
	}

	userRoles, err := s.userRoles.FindByUserID(ctx, user.ID)
	if err != nil {
		return input.GetUserResult{}, fmt.Errorf("find user roles: %w", err)
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roleList, err := s.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return input.GetUserResult{}, fmt.Errorf("find roles: %w", err)
	}

	item, err := domainservice.BuildUserItem(*user, userRoles, roleList)
	if err != nil {
		return input.GetUserResult{}, err
	}

	return input.GetUserResult{
		ID:       item.ID,
		Email:    item.Email,
		IsActive: item.IsActive,
		Roles:    item.Roles,
	}, nil
}
