package auth

import (
	"context"
	"fmt"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

type MeUseCase struct {
	users    output.UserRepository
	resolver *common.PermissionResolver
}

func NewMeUseCase(
	users output.UserRepository,
	resolver *common.PermissionResolver,
) *MeUseCase {
	return &MeUseCase{
		users:    users,
		resolver: resolver,
	}
}

func (uc *MeUseCase) Me(ctx context.Context, claims value.AccessClaims) (dto.CurrentUser, error) {
	user, err := uc.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return dto.CurrentUser{}, fmt.Errorf("find user: %w", err)
	}

	if !user.IsActive {
		return dto.CurrentUser{}, domain.ErrUserInactive
	}

	roles, perms, err := uc.resolver.Resolve(ctx, user.ID)
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
