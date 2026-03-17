package auth

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *AuthService) Logout(ctx context.Context, in input.LogoutInput) error {
	hash := s.tokenH.Hash(in.RawRefreshToken)

	session, err := s.sessions.FindByTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			// Идемпотентность: сессия уже не существует — logout считается успешным.
			return nil
		}
		return fmt.Errorf("find refresh session: %w", err)
	}

	if session.IsRevoked() {
		return nil
	}

	if err := s.sessions.Revoke(ctx, session.ID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
