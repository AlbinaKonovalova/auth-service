package auth

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

type LogoutUseCase struct {
	sessions output.RefreshSessionRepository
	hasher   output.TokenHasher
}

func NewLogoutUseCase(
	sessions output.RefreshSessionRepository,
	hasher output.TokenHasher,
) *LogoutUseCase {
	return &LogoutUseCase{sessions: sessions, hasher: hasher}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, in input.LogoutInput) error {
	hash := uc.hasher.Hash(in.RawRefreshToken)

	session, err := uc.sessions.FindByTokenHash(ctx, hash)
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

	if err := uc.sessions.Revoke(ctx, session.ID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
