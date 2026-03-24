package password_reset

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, in input.RequestPasswordResetInput) error {
	now := s.clock.Now()

	email, err := value.NewEmail(in.Email)
	if err != nil {
		return nil
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return fmt.Errorf("find user by email: %w", err)
	}

	rawToken, tokenHash, err := s.tokenProvider.GeneratePasswordResetToken(ctx)
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}

	tokenID := s.uuidGen.New()

	if err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		if _, err := s.userRepo.FindByIDForUpdate(ctx, user.ID); err != nil {
			return fmt.Errorf("lock user for reset: %w", err)
		}

		if err := s.resetRepo.InvalidateByUserID(ctx, user.ID, now); err != nil {
			return fmt.Errorf("invalidate old reset tokens: %w", err)
		}

		resetToken, err := entity.NewPasswordResetToken(tokenID, user.ID, tokenHash.String(), now, s.tokenTTL)
		if err != nil {
			return fmt.Errorf("create reset token entity: %w", err)
		}

		if err := s.resetRepo.Create(ctx, resetToken); err != nil {
			return fmt.Errorf("save reset token: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	if err := s.mailer.SendPasswordResetEmail(ctx, user.Email, rawToken); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return nil
}
