package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeactivateUser деактивирует пользователя и отзывает все его активные refresh sessions.
// Оба действия выполняются в одной транзакции.
//
// Идемпотентность:
//   - если пользователь уже неактивен — Deactivate не вызывается
//   - RevokeAllByUserID вызывается всегда, чтобы гарантировать отсутствие живых сессий
//     даже если данные оказались в рассогласованном состоянии
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	return s.tx.RunInTx(ctx, func(ctx context.Context) error {
		user, err := s.users.FindByIDForUpdate(ctx, id)
		if err != nil {
			return fmt.Errorf("deactivate user: %w", err)
		}

		if user.IsActive {
			if err := s.users.Deactivate(ctx, id); err != nil {
				return fmt.Errorf("deactivate user: %w", err)
			}
		}

		if err := s.sessions.RevokeAllByUserID(ctx, id); err != nil {
			return fmt.Errorf("revoke sessions on deactivate: %w", err)
		}

		return nil
	})
}
