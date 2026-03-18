package service

import (
	"sort"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/google/uuid"
)

type UserListItem struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

type UserList struct {
	Items   []UserListItem
	Total   int
	Page    int
	PerPage int
}

func BuildUserList(
	users []entity.User,
	userRoles []entity.UserRole,
	roles []entity.Role,
	total int,
	page int,
	perPage int,
) (UserList, error) {
	if total < 0 {
		return UserList{}, domain.ErrInvalidUserList
	}

	if len(users) == 0 {
		return UserList{
			Items:   []UserListItem{},
			Total:   total,
			Page:    page,
			PerPage: perPage,
		}, nil
	}

	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, r := range roles {
		roleByID[r.ID] = r
	}

	rolesByUser := make(map[uuid.UUID][]string)
	for _, ur := range userRoles {
		role, ok := roleByID[ur.RoleID]
		if !ok {
			return UserList{}, domain.ErrRoleNotFound
		}
		rolesByUser[ur.UserID] = append(rolesByUser[ur.UserID], role.Code)
	}

	items := make([]UserListItem, len(users))
	for i, u := range users {
		roleCodes := rolesByUser[u.ID]
		if roleCodes == nil {
			roleCodes = []string{}
		}

		sort.Strings(roleCodes)

		items[i] = UserListItem{
			ID:       u.ID,
			Email:    u.Email,
			IsActive: u.IsActive,
			Roles:    roleCodes,
		}
	}

	return UserList{
		Items:   items,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// BuildUserItem собирает доменное представление одного пользователя с его ролями.
// Используется в GetUser сценарии.
func BuildUserItem(
	user entity.User,
	userRoles []entity.UserRole,
	roles []entity.Role,
) (UserListItem, error) {
	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, r := range roles {
		roleByID[r.ID] = r
	}

	roleCodes := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		role, ok := roleByID[ur.RoleID]
		if !ok {
			return UserListItem{}, domain.ErrRoleNotFound
		}
		roleCodes = append(roleCodes, role.Code)
	}

	sort.Strings(roleCodes)

	return UserListItem{
		ID:       user.ID,
		Email:    user.Email,
		IsActive: user.IsActive,
		Roles:    roleCodes,
	}, nil
}

func ValidateUserListFilters(
	filters dto.UserListFilters,
	foundRoles []entity.Role,
) (dto.UserListFilters, error) {
	if filters.Role != "" {
		roleCode, err := value.NewRoleCode(filters.Role)
		if err != nil {
			return dto.UserListFilters{}, err
		}
		filters.Role = roleCode.String()

		if len(foundRoles) == 0 {
			return dto.UserListFilters{}, domain.ErrRoleNotFound
		}
	}

	filters.Page, filters.PerPage = NormalizeUserListPagination(filters.Page, filters.PerPage)

	return filters, nil
}

const (
	MinUserListPage        = 1
	MinUserListPerPage     = 1
	DefaultUserListPage    = 1
	DefaultUserListPerPage = 20
	MaxUserListPerPage     = 100
)

func NormalizeUserListPagination(page, perPage int) (int, int) {
	if page < MinUserListPage {
		page = DefaultUserListPage
	}

	if perPage < MinUserListPerPage {
		perPage = DefaultUserListPerPage
	}

	if perPage > MaxUserListPerPage {
		perPage = MaxUserListPerPage
	}

	return page, perPage
}
