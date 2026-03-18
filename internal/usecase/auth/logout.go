package auth

import (
	"context"
	"fmt"

	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *AuthService) Logout(ctx context.Context, in input.LogoutInput) error {
	hash := s.tokenH.Hash(in.RawRefreshToken)

	session, err := s.sessions.FindByTokenHash(ctx, hash)

	canRevoke, err := domainservice.CanLogout(err, session)
	if err != nil {
		return fmt.Errorf("find refresh session: %w", err)
	}
	if !canRevoke {
		return nil
	}

	if err := s.sessions.Revoke(ctx, session.ID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
