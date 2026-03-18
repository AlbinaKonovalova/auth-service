package auth

import (
	"context"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func (s *AuthService) Me(ctx context.Context, claims value.AccessClaims) (dto.CurrentUser, error) {
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return dto.CurrentUser{}, fmt.Errorf("find user: %w", err)
	}

	if err := user.EnsureActive(); err != nil {
		return dto.CurrentUser{}, err
	}

	roles, perms, err := s.resolver.Resolve(ctx, user.ID)
	if err != nil {
		return dto.CurrentUser{}, fmt.Errorf("resolve permissions: %w", err)
	}

	return dto.CurrentUser{
		ID:          user.ID,
		Email:       user.Email,
		IsActive:    user.IsActive,
		Roles:       roles,
		Permissions: perms,
	}, nil
}
