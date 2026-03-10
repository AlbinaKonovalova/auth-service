package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

const refreshTTL = 7 * 24 * time.Hour

type LoginUseCase struct {
	users    output.UserRepository
	sessions output.RefreshSessionRepository
	hasher   output.PasswordHasher
	tokens   output.TokenProvider
	tx       output.TxManager
	clock    output.Clock
	uuid     output.UUIDGenerator
	resolver *common.PermissionResolver
}

func NewLoginUseCase(
	users output.UserRepository,
	sessions output.RefreshSessionRepository,
	hasher output.PasswordHasher,
	tokens output.TokenProvider,
	tx output.TxManager,
	clock output.Clock,
	uuid output.UUIDGenerator,
	resolver *common.PermissionResolver,
) *LoginUseCase {
	return &LoginUseCase{
		users:    users,
		sessions: sessions,
		hasher:   hasher,
		tokens:   tokens,
		tx:       tx,
		clock:    clock,
		uuid:     uuid,
		resolver: resolver,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, in input.LoginInput) (dto.LoginResult, string, error) {
	email, err := value.NewEmail(in.Email)
	if err != nil {
		return dto.LoginResult{}, "", domain.ErrInvalidCredentials
	}

	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return dto.LoginResult{}, "", domain.ErrInvalidCredentials
		}
		return dto.LoginResult{}, "", fmt.Errorf("find user: %w", err)
	}

	if !user.IsActive {
		return dto.LoginResult{}, "", domain.ErrUserInactive
	}

	ok, err := uc.hasher.Verify(in.Password, user.PasswordHash)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("verify password: %w", err)
	}
	if !ok {
		return dto.LoginResult{}, "", domain.ErrInvalidCredentials
	}

	roles, perms, err := uc.resolver.Resolve(ctx, user.ID)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("resolve permissions: %w", err)
	}

	claims := value.AccessClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Roles:       roles,
		Permissions: perms,
	}

	accessToken, expiresIn, err := uc.tokens.GenerateAccessToken(ctx, claims)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("generate access token: %w", err)
	}

	var result dto.LoginResult
	var rawRefresh string

	err = uc.tx.RunInTx(ctx, func(ctx context.Context) error {
		if err := uc.sessions.DeleteExpiredAndRevoked(ctx, user.ID); err != nil {
			return fmt.Errorf("cleanup sessions: %w", err)
		}

		raw, hash, err := uc.tokens.GenerateRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		now := uc.clock.Now()
		session := entity.RefreshSession{
			ID:        uc.uuid.New(),
			UserID:    user.ID,
			TokenHash: hash,
			ExpiresAt: now.Add(refreshTTL),
			CreatedAt: now,
		}

		if err := uc.sessions.Save(ctx, session); err != nil {
			return fmt.Errorf("save refresh session: %w", err)
		}

		rawRefresh = raw
		result = dto.LoginResult{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
			User: dto.CurrentUser{
				ID:          user.ID,
				Email:       user.Email,
				IsActive:    user.IsActive,
				Roles:       roles,
				Permissions: perms,
			},
		}

		return nil
	})
	if err != nil {
		return dto.LoginResult{}, "", err
	}

	return result, rawRefresh, nil
}
