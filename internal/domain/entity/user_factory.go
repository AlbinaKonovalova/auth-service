package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type NewUserAggregate struct {
	User      User
	UserRoles []UserRole
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

	if _, err := value.NewPassword(rawPassword); err != nil {
		return NewUserAggregate{}, err
	}

	requestedRoleCodes, err := NormalizeRequestedRoleCodes(rawRoleCodes)
	if err != nil {
		return NewUserAggregate{}, err
	}

	if err := ensureUserHasAtLeastOneRole(requestedRoleCodes); err != nil {
		return NewUserAggregate{}, err
	}

	if err := EnsureAllRequestedRolesExist(requestedRoleCodes, foundRoles); err != nil {
		return NewUserAggregate{}, err
	}

	user, err := NewUser(userID, email, passwordHash, now)
	if err != nil {
		return NewUserAggregate{}, err
	}

	assignments := make([]UserRole, 0, len(foundRoles))
	for _, role := range foundRoles {
		userRole, err := NewUserRole(userID, role.ID, now)
		if err != nil {
			return NewUserAggregate{}, err
		}
		assignments = append(assignments, userRole)
	}

	return NewUserAggregate{
		User:      user,
		UserRoles: assignments,
	}, nil
}

func ensureUserHasAtLeastOneRole(requested []value.RoleCode) error {
	if len(requested) == 0 {
		return domain.ErrUserMustHaveRole
	}

	return nil
}
