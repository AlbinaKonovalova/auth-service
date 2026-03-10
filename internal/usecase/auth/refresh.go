package auth

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

type RefreshUseCase struct {
	users    output.UserRepository
	sessions output.RefreshSessionRepository
	tokens   output.TokenProvider
	hasher   output.TokenHasher
	tx       output.TxManager
	clock    output.Clock
	uuid     output.UUIDGenerator
	resolver *common.PermissionResolver
}

func NewRefreshUseCase(
	users output.UserRepository,
	sessions output.RefreshSessionRepository,
	tokens output.TokenProvider,
	hasher output.TokenHasher,
	tx output.TxManager,
	clock output.Clock,
	uuid output.UUIDGenerator,
	resolver *common.PermissionResolver,
) *RefreshUseCase {
	return &RefreshUseCase{
		users:    users,
		sessions: sessions,
		tokens:   tokens,
		hasher:   hasher,
		tx:       tx,
		clock:    clock,
		uuid:     uuid,
		resolver: resolver,
	}
}

func (uc *RefreshUseCase) Refresh(ctx context.Context, in input.RefreshInput) (dto.RefreshResult, string, error) {
	hash := uc.hasher.Hash(in.RawRefreshToken)

	var result dto.RefreshResult
	var rawRefresh string

	err := uc.tx.RunInTx(ctx, func(ctx context.Context) error {
		session, err := uc.sessions.FindByTokenHashForUpdate(ctx, hash)
		if err != nil {
			if errors.Is(err, domain.ErrRefreshTokenNotFound) {
				return domain.ErrRefreshTokenNotFound
			}
			return fmt.Errorf("find refresh session: %w", err)
		}

		now := uc.clock.Now()

		if session.IsRevoked() {
			return domain.ErrRefreshTokenRevoked
		}
		if session.IsExpired(now) {
			return domain.ErrRefreshTokenExpired
		}

		user, err := uc.users.FindByID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}
		if !user.IsActive {
			return domain.ErrUserInactive
		}

		roles, perms, err := uc.resolver.Resolve(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("resolve permissions: %w", err)
		}

		claims := value.AccessClaims{
			UserID:      user.ID,
			Email:       user.Email,
			Roles:       roles,
			Permissions: perms,
		}

		accessToken, expiresIn, err := uc.tokens.GenerateAccessToken(ctx, claims)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		if err := uc.sessions.Revoke(ctx, session.ID); err != nil {
			return fmt.Errorf("revoke old session: %w", err)
		}

		raw, newHash, err := uc.tokens.GenerateRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		newSession := entity.RefreshSession{
			ID:        uc.uuid.New(),
			UserID:    user.ID,
			TokenHash: newHash,
			ExpiresAt: now.Add(refreshTTL),
			CreatedAt: now,
		}

		if err := uc.sessions.Save(ctx, newSession); err != nil {
			return fmt.Errorf("save new session: %w", err)
		}

		rawRefresh = raw
		result = dto.RefreshResult{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
		}

		return nil
	})
	if err != nil {
		return dto.RefreshResult{}, "", err
	}

	return result, rawRefresh, nil
}
