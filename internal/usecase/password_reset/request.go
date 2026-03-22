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

// RequestPasswordReset инициирует сброс пароля.
//
// Сценарий:
//  1. Нормализовать email через value object — невалидный формат молча возвращает success
//  2. Найти пользователя по email
//     — ErrUserNotFound → return nil (no user enumeration)
//     — любая другая ошибка repo → вернуть ошибку наверх (DB failure не скрывается)
//  3. Сгенерировать raw token и hash вне транзакции
//  4. Внутри транзакции с блокировкой строки пользователя (FOR UPDATE):
//     а) перечитать пользователя под lock — используется для сериализации
//     параллельных reset request для одного user
//     б) инвалидировать все активные reset tokens пользователя
//     в) создать новый reset token через доменный конструктор
//     г) сохранить token в репозитории
//  5. Отправить email с raw token после успешной транзакции
//
// Метод возвращает nil когда пользователь не найден (no user enumeration).
// Ошибки инфраструктуры (DB, SMTP) возвращаются как есть.
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, in input.RequestPasswordResetInput) error {
	now := s.clock.Now()

	// Невалидный email — молча успех (no user enumeration)
	email, err := value.NewEmail(in.Email)
	if err != nil {
		return nil
	}

	// Поиск пользователя по email
	// ErrUserNotFound → молча успех (no user enumeration)
	// любая другая ошибка → честная internal ошибка, не маскируется
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return fmt.Errorf("find user by email: %w", err)
	}

	// Генерируем raw token и hash вне транзакции
	rawToken, tokenHash, err := s.tokenProvider.GeneratePasswordResetToken(ctx)
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}

	tokenID := s.uuidGen.New()

	// Транзакция с блокировкой строки пользователя:
	// FindByIDForUpdate используется для сериализации параллельных reset request
	// для одного user — параллельные транзакции будут ждать снятия блокировки.
	if err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		// Lock user row для сериализации параллельных request
		if _, err := s.userRepo.FindByIDForUpdate(ctx, user.ID); err != nil {
			return fmt.Errorf("lock user for reset: %w", err)
		}

		// Инвалидировать старые токены
		if err := s.resetRepo.InvalidateByUserID(ctx, user.ID, now); err != nil {
			return fmt.Errorf("invalidate old reset tokens: %w", err)
		}

		// Доменная сущность создаётся через конструктор с валидацией
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

	// Отправка письма после успешной транзакции
	// Ошибка отправки — честная internal error, не маскируется
	if err := s.mailer.SendPasswordResetEmail(ctx, user.Email, rawToken); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return nil
}
