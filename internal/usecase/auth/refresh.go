package auth

import (
	"context"
	"fmt"

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
			return fmt.Errorf("find refresh session: %w", err)
		}

		now := s.clock.Now()

		if err := session.EnsureUsable(now); err != nil {
			return err
		}

		user, err := s.users.FindByID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		if err := user.EnsureActive(); err != nil {
			return err
		}

		roles, perms, err := s.resolver.Resolve(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("resolve permissions: %w", err)
		}

		claims, err := value.NewAccessClaims(user.ID, user.Email, roles, perms)
		if err != nil {
			return err
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

		newSession, err := entity.NewRefreshSession(
			s.uuid.New(),
			user.ID,
			newHash,
			now,
			s.refreshTTL,
		)
		if err != nil {
			return err
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
