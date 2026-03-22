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
// Повторный confirm тем же token вернёт ErrResetTokenUsed.
func (s *PasswordResetService) ConfirmPasswordReset(ctx context.Context, in input.ConfirmPasswordResetInput) error {
	now := s.clock.Now()

	// Валидация нового пароля через value object
	if _, err := value.NewPassword(in.NewPassword); err != nil {
		return err
	}

	// Hash raw token — в repo уходит только hash, raw token нигде не сохраняется
	tokenHash := s.tokenHasher.Hash(in.Token)

	return s.tx.RunInTx(ctx, func(ctx context.Context) error {
		// Найти и заблокировать reset token по hash (FOR UPDATE).
		// Lock берётся именно на token row — token определяет одноразовость сценария.
		// Параллельные confirm будут ждать снятия блокировки.
		resetToken, err := s.resetRepo.FindByTokenHashForUpdate(ctx, tokenHash.String())
		if err != nil {
			if errors.Is(err, domain.ErrResetTokenNotFound) {
				return domain.ErrResetTokenNotFound
			}
			return fmt.Errorf("find reset token for update: %w", err)
		}

		// Проверить применимость через доменный метод (expired → used, в таком порядке)
		if err := resetToken.EnsureUsable(now); err != nil {
			return err
		}

		// Посчитать новый password hash через hasher
		newPasswordHash, err := s.hasher.Hash(in.NewPassword)
		if err != nil {
			return fmt.Errorf("hash new password: %w", err)
		}

		// Обновить password hash пользователя
		if err := s.userRepo.UpdatePasswordHash(ctx, resetToken.UserID, newPasswordHash); err != nil {
			return fmt.Errorf("update password hash: %w", err)
		}

		// Пометить reset token как used — одноразовость гарантируется lock + RowsAffected
		if err := s.resetRepo.MarkUsed(ctx, resetToken.ID, now); err != nil {
			return fmt.Errorf("mark reset token used: %w", err)
		}

		// Отозвать все refresh sessions пользователя
		if err := s.sessionRepo.RevokeAllByUserID(ctx, resetToken.UserID); err != nil {
			return fmt.Errorf("revoke refresh sessions: %w", err)
		}

		return nil
	})
}
