package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (domainservice.AdminUserView, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return domainservice.AdminUserView{}, fmt.Errorf("find user by id: %w", err)
	}

	userRoles, err := s.userRoles.FindByUserID(ctx, user.ID)
	if err != nil {
		return domainservice.AdminUserView{}, fmt.Errorf("find user roles: %w", err)
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roleList, err := s.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return domainservice.AdminUserView{}, fmt.Errorf("find roles: %w", err)
	}

	return domainservice.BuildAdminUserView(*user, userRoles, roleList)
}
