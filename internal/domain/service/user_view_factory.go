package service

import (
	"sort"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/google/uuid"
)

type AdminUserView struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

type UserList struct {
	Items   []AdminUserView
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
			Items:   []AdminUserView{},
			Total:   total,
			Page:    page,
			PerPage: perPage,
		}, nil
	}

	roleByID := buildRoleByID(roles)

	rolesByUser := make(map[uuid.UUID][]entity.UserRole)
	for _, ur := range userRoles {
		rolesByUser[ur.UserID] = append(rolesByUser[ur.UserID], ur)
	}

	items := make([]AdminUserView, len(users))
	for i, u := range users {
		roleCodes, err := buildSortedUniqueRoleCodes(rolesByUser[u.ID], roleByID)
		if err != nil {
			return UserList{}, err
		}

		items[i] = AdminUserView{
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

// BuildAdminUserView собирает доменное представление одного пользователя с его ролями.
// Используется в GetUser и CreateUser сценариях.
func BuildAdminUserView(
	user entity.User,
	userRoles []entity.UserRole,
	roles []entity.Role,
) (AdminUserView, error) {
	roleByID := buildRoleByID(roles)

	roleCodes, err := buildSortedUniqueRoleCodes(userRoles, roleByID)
	if err != nil {
		return AdminUserView{}, err
	}

	return AdminUserView{
		ID:       user.ID,
		Email:    user.Email,
		IsActive: user.IsActive,
		Roles:    roleCodes,
	}, nil
}

func buildRoleByID(roles []entity.Role) map[uuid.UUID]entity.Role {
	roleByID := make(map[uuid.UUID]entity.Role, len(roles))
	for _, r := range roles {
		roleByID[r.ID] = r
	}
	return roleByID
}

func buildSortedUniqueRoleCodes(
	userRoles []entity.UserRole,
	roleByID map[uuid.UUID]entity.Role,
) ([]string, error) {
	if len(userRoles) == 0 {
		return []string{}, nil
	}

	seen := make(map[string]struct{}, len(userRoles))
	roleCodes := make([]string, 0, len(userRoles))

	for _, ur := range userRoles {
		role, ok := roleByID[ur.RoleID]
		if !ok {
			return nil, domain.ErrDataIntegrityViolation
		}

		if _, exists := seen[role.Code]; exists {
			continue
		}
		seen[role.Code] = struct{}{}
		roleCodes = append(roleCodes, role.Code)
	}

	sort.Strings(roleCodes)

	return roleCodes, nil
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
