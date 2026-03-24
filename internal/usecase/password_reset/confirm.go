package password_reset

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// ConfirmPasswordReset подтверждает сброс пароля.
//
// Сценарий (read-check-write внутри одной транзакции):
//  1. Валидировать новый пароль через value object
//  2. Захешировать raw token через token hasher
//  3. Внутри транзакции:
//     а) найти и заблокировать reset token по hash (FOR UPDATE) —
//     lock берётся именно на token row, так как token определяет одноразовость сценария
//     б) проверить применимость через domain entity (expired → used, в таком порядке)
//     в) посчитать новый password hash через hasher
//     г) обновить password hash пользователя
//     д) пометить reset token как used — MarkUsed дополнительно защищает от race через RowsAffected
//     е) отозвать все refresh sessions пользователя
//

func (s *PasswordResetService) ConfirmPasswordReset(ctx context.Context, in input.ConfirmPasswordResetInput) error {
	now := s.clock.Now()

	if _, err := value.NewPassword(in.NewPassword); err != nil {
		return err
	}

	// Hash raw token — в repo уходит только hash, raw token нигде не сохраняется
	tokenHash := s.tokenHasher.Hash(in.Token)

	return s.tx.RunInTx(ctx, func(ctx context.Context) error {

		resetToken, err := s.resetRepo.FindByTokenHashForUpdate(ctx, tokenHash.String())
		if err != nil {
			if errors.Is(err, domain.ErrResetTokenNotFound) {
				return domain.ErrResetTokenNotFound
			}
			return fmt.Errorf("find reset token for update: %w", err)
		}

		if err := resetToken.EnsureUsable(now); err != nil {
			return err
		}

		newPasswordHash, err := s.hasher.Hash(in.NewPassword)
		if err != nil {
			return fmt.Errorf("hash new password: %w", err)
		}

		if err := s.userRepo.UpdatePasswordHash(ctx, resetToken.UserID, newPasswordHash); err != nil {
			return fmt.Errorf("update password hash: %w", err)
		}

		if err := s.resetRepo.MarkUsed(ctx, resetToken.ID, now); err != nil {
			return fmt.Errorf("mark reset token used: %w", err)
		}

		if err := s.sessionRepo.RevokeAllByUserID(ctx, resetToken.UserID); err != nil {
			return fmt.Errorf("revoke refresh sessions: %w", err)
		}

		return nil
	})
}
