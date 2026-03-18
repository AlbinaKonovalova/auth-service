package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *AuthService) Login(ctx context.Context, in input.LoginInput) (dto.LoginResult, string, error) {
	email, err := value.NewEmail(in.Email)
	if err != nil {
		return dto.LoginResult{}, "", domain.ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		err = domainservice.MapLoginLookupError(err)
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return dto.LoginResult{}, "", err
		}
		return dto.LoginResult{}, "", fmt.Errorf("find user: %w", err)
	}

	if err := user.EnsureActive(); err != nil {
		return dto.LoginResult{}, "", err
	}

	ok, err := s.hasher.Verify(in.Password, user.PasswordHash)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("verify password: %w", err)
	}
	if !ok {
		return dto.LoginResult{}, "", domain.ErrInvalidCredentials
	}

	roles, perms, err := s.resolver.Resolve(ctx, user.ID)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("resolve permissions: %w", err)
	}

	claims, err := value.NewAccessClaims(user.ID, user.Email, roles, perms)
	if err != nil {
		return dto.LoginResult{}, "", err
	}

	accessToken, _, err := s.tokens.GenerateAccessToken(ctx, claims)
	if err != nil {
		return dto.LoginResult{}, "", fmt.Errorf("generate access token: %w", err)
	}

	var result dto.LoginResult
	var rawRefresh string

	err = s.tx.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.sessions.DeleteExpiredAndRevoked(ctx, user.ID); err != nil {
			return fmt.Errorf("cleanup sessions: %w", err)
		}

		raw, hash, err := s.tokens.GenerateRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		now := s.clock.Now()
		session, err := entity.NewRefreshSession(
			s.uuid.New(),
			user.ID,
			hash,
			now,
			s.refreshTTL,
		)
		if err != nil {
			return err
		}

		if err := s.sessions.Save(ctx, session); err != nil {
			return fmt.Errorf("save refresh session: %w", err)
		}

		rawRefresh = raw
		result = dto.LoginResult{
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
		return dto.LoginResult{}, "", err
	}

	return result, rawRefresh, nil
}
