package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ActivateUser активирует пользователя.
// Сценарий идемпотентен: если пользователь уже активен — возвращает success.
func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	return s.tx.RunInTx(ctx, func(ctx context.Context) error {
		user, err := s.users.FindByIDForUpdate(ctx, id)
		if err != nil {
			return fmt.Errorf("activate user: %w", err)
		}

		if user.IsActive {
			return nil
		}

		if err := s.users.Activate(ctx, id); err != nil {
			return fmt.Errorf("activate user: %w", err)
		}

		return nil
	})
}
