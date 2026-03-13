package entity

import (
	"time"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/google/uuid"
)

type NewUserAggregate struct {
	User          User
	UserRoles     []UserRole
	AssignedRoles []string
}

func BuildNewUserAggregate(
	userID uuid.UUID,
	rawEmail string,
	rawPassword string,
	passwordHash string,
	rawRoleCodes []string,
	foundRoles []Role,
	now time.Time,
) (NewUserAggregate, error) {
	email, err := value.NewEmail(rawEmail)
	if err != nil {
		return NewUserAggregate{}, err
	}

	password, err := value.NewPassword(rawPassword)
	if err != nil {
		return NewUserAggregate{}, err
	}
	_ = password // raw password validated intentionally

	requestedRoleCodes, err := normalizeRoleCodes(rawRoleCodes)
	if err != nil {
		return NewUserAggregate{}, err
	}

	if err := ensureAllRequestedRolesExist(requestedRoleCodes, foundRoles); err != nil {
		return NewUserAggregate{}, err
	}

	user, err := NewUser(userID, email, passwordHash, now)
	if err != nil {
		return NewUserAggregate{}, err
	}

	assignments := make([]UserRole, 0, len(foundRoles))
	assignedCodes := make([]string, 0, len(foundRoles))

	for _, role := range foundRoles {
		userRole, err := NewUserRole(userID, role.ID, now)
		if err != nil {
			return NewUserAggregate{}, err
		}

		assignments = append(assignments, userRole)
		assignedCodes = append(assignedCodes, role.Code)
	}

	return NewUserAggregate{
		User:          user,
		UserRoles:     assignments,
		AssignedRoles: assignedCodes,
	}, nil
}

func normalizeRoleCodes(raw []string) ([]value.RoleCode, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	seen := make(map[string]struct{}, len(raw))
	result := make([]value.RoleCode, 0, len(raw))

	for _, item := range raw {
		code, err := value.NewRoleCode(item)
		if err != nil {
			return nil, err
		}

		if _, exists := seen[code.String()]; exists {
			return nil, domain.ErrDuplicateRoleCode
		}

		seen[code.String()] = struct{}{}
		result = append(result, code)
	}

	return result, nil
}

func ensureAllRequestedRolesExist(requested []value.RoleCode, found []Role) error {
	if len(requested) == 0 {
		return nil
	}

	foundSet := make(map[string]struct{}, len(found))
	for _, role := range found {
		foundSet[role.Code] = struct{}{}
	}

	for _, code := range requested {
		if _, ok := foundSet[code.String()]; !ok {
			return domain.ErrRoleNotFound
		}
	}

	return nil
}
