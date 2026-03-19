package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

func (s *UserService) ListUsers(ctx context.Context, filters dto.UserListFilters) (domainservice.UserList, error) {
	var foundRolesByFilter []entity.Role
	var err error

	if filters.Role != "" {
		foundRolesByFilter, err = s.roles.FindByCodes(ctx, []string{filters.Role})
		if err != nil {
			return domainservice.UserList{}, fmt.Errorf("find role by code: %w", err)
		}
	}

	filters, err = domainservice.ValidateUserListFilters(filters, foundRolesByFilter)
	if err != nil {
		return domainservice.UserList{}, err
	}

	foundUsers, total, err := s.users.List(ctx, filters)
	if err != nil {
		return domainservice.UserList{}, fmt.Errorf("list users: %w", err)
	}

	userIDs := make([]uuid.UUID, len(foundUsers))
	for i, u := range foundUsers {
		userIDs[i] = u.ID
	}

	userRoleRows, err := s.userRoles.FindByUserIDs(ctx, userIDs)
	if err != nil {
		return domainservice.UserList{}, fmt.Errorf("find user roles: %w", err)
	}

	roleIDSet := make(map[uuid.UUID]struct{}, len(userRoleRows))
	for _, ur := range userRoleRows {
		roleIDSet[ur.RoleID] = struct{}{}
	}

	roleIDs := make([]uuid.UUID, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}

	roleList, err := s.roles.FindByIDs(ctx, roleIDs)
	if err != nil {
		return domainservice.UserList{}, fmt.Errorf("find roles: %w", err)
	}

	return domainservice.BuildUserList(
		foundUsers,
		userRoleRows,
		roleList,
		total,
		filters.Page,
		filters.PerPage,
	)
}
