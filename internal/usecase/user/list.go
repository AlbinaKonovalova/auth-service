package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *UserService) ListUsers(ctx context.Context, filters dto.UserListFilters) (input.ListUsersResult, error) {
	var foundRolesByFilter []entity.Role
	var err error

	if filters.Role != "" {
		foundRolesByFilter, err = s.roles.FindByCodes(ctx, []string{filters.Role})
		if err != nil {
			return input.ListUsersResult{}, fmt.Errorf("find role by code: %w", err)
		}
	}

	filters, err = domainservice.ValidateUserListFilters(filters, foundRolesByFilter)
	if err != nil {
		return input.ListUsersResult{}, err
	}

	foundUsers, total, err := s.users.List(ctx, filters)
	if err != nil {
		return input.ListUsersResult{}, fmt.Errorf("list users: %w", err)
	}

	userIDs := make([]uuid.UUID, len(foundUsers))
	for i, u := range foundUsers {
		userIDs[i] = u.ID
	}

	userRoleRows, err := s.userRoles.FindByUserIDs(ctx, userIDs)
	if err != nil {
		return input.ListUsersResult{}, fmt.Errorf("find user roles: %w", err)
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
		return input.ListUsersResult{}, fmt.Errorf("find roles: %w", err)
	}

	userList, err := domainservice.BuildUserList(
		foundUsers,
		userRoleRows,
		roleList,
		total,
		filters.Page,
		filters.PerPage,
	)
	if err != nil {
		return input.ListUsersResult{}, err
	}

	results := make([]input.GetUserResult, len(userList.Items))
	for i, item := range userList.Items {
		results[i] = input.GetUserResult{
			ID:       item.ID,
			Email:    item.Email,
			IsActive: item.IsActive,
			Roles:    item.Roles,
		}
	}

	return input.ListUsersResult{
		Users:   results,
		Total:   userList.Total,
		Page:    userList.Page,
		PerPage: userList.PerPage,
	}, nil
}
