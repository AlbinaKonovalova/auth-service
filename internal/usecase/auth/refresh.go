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
)

func (s *AuthService) Refresh(ctx context.Context, in input.RefreshInput) (dto.RefreshResult, string, error) {
	hash := s.tokenH.Hash(in.RawRefreshToken)

	var result dto.RefreshResult
	var rawRefresh string

	err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		session, err := s.sessions.FindByTokenHashForUpdate(ctx, hash)
		if err != nil {
			if errors.Is(err, domain.ErrRefreshTokenNotFound) {
				return domain.ErrRefreshTokenNotFound
			}
			return fmt.Errorf("find refresh session: %w", err)
		}

		now := s.clock.Now()

		if session.IsRevoked() {
			return domain.ErrRefreshTokenRevoked
		}
		if session.IsExpired(now) {
			return domain.ErrRefreshTokenExpired
		}

		user, err := s.users.FindByID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}
		if !user.IsActive {
			return domain.ErrUserInactive
		}

		roles, perms, err := s.resolver.Resolve(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("resolve permissions: %w", err)
		}

		claims := value.AccessClaims{
			UserID:      user.ID,
			Email:       user.Email,
			Roles:       roles,
			Permissions: perms,
		}

		accessToken, _, err := s.tokens.GenerateAccessToken(ctx, claims)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		if err := s.sessions.Revoke(ctx, session.ID); err != nil {
			return fmt.Errorf("revoke old session: %w", err)
		}

		raw, newHash, err := s.tokens.GenerateRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		newSession := entity.RefreshSession{
			ID:        s.uuid.New(),
			UserID:    user.ID,
			TokenHash: newHash,
			ExpiresAt: now.Add(refreshTTL),
			CreatedAt: now,
		}

		if err := s.sessions.Save(ctx, newSession); err != nil {
			return fmt.Errorf("save new session: %w", err)
		}

		rawRefresh = raw
		result = dto.RefreshResult{
			AccessToken: accessToken,
			User: dto.CurrentUser{
				ID:          user.ID,
				Roles:       roles,
				Permissions: perms,
			},
		}

		return nil
	})
	if err != nil {
		return dto.RefreshResult{}, "", err
	}

	return result, rawRefresh, nil
}
