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
		return input.CreateUserResult{}, fmt.Errorf("hash password: %w", err)
	}

	now := s.clock.Now()
	userID := s.uuid.New()

	foundRoles, err := s.roles.FindByCodes(ctx, in.Roles)
	if err != nil {
		return input.CreateUserResult{}, fmt.Errorf("find roles by codes: %w", err)
	}

	aggregate, err := entity.BuildNewUserAggregate(
		userID,
		in.Email,
		in.Password,
		hash,
		in.Roles,
		foundRoles,
		now,
	)
	if err != nil {
		return input.CreateUserResult{}, err
	}

	err = s.tx.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.users.Create(txCtx, aggregate.User); err != nil {
			return err
		}

		for _, userRole := range aggregate.UserRoles {
			if err := s.userRoles.Assign(txCtx, userRole); err != nil {
				return fmt.Errorf("assign user role: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return input.CreateUserResult{}, err
	}

	return input.CreateUserResult{
		ID:       aggregate.User.ID,
		Email:    aggregate.User.Email,
		IsActive: aggregate.User.IsActive,
		Roles:    aggregate.AssignedRoles,
	}, nil
}
