package input

import "context"

// RequestPasswordResetInput — входные данные для запроса сброса пароля.
type RequestPasswordResetInput struct {
	Email string
}

// ConfirmPasswordResetInput — входные данные для подтверждения сброса пароля.
type ConfirmPasswordResetInput struct {
	// Token — сырой токен из ссылки в письме (не хеш).
	Token       string
	NewPassword string
}

// PasswordResetUseCase — входной контракт для password reset flow.
type PasswordResetUseCase interface {
	// RequestPasswordReset инициирует сброс пароля.
	// На уровне foundation (commit 13) это контракт сценария.
	// Реализация (создание токена, отправка письма) идёт в commit 14.
	RequestPasswordReset(ctx context.Context, in RequestPasswordResetInput) error

	// ConfirmPasswordReset подтверждает сброс пароля.
	// Проверяет токен, меняет пароль, помечает токен использованным
	// и отзывает все refresh sessions пользователя.
	// Возвращает domain.ErrResetTokenNotFound если токен не найден.
	// Возвращает domain.ErrResetTokenExpired если токен истёк.
	// Возвращает domain.ErrResetTokenUsed если токен уже был использован.
	// Возвращает domain.ErrInvalidPassword если новый пароль не прошёл валидацию.
	ConfirmPasswordReset(ctx context.Context, in ConfirmPasswordResetInput) error
}
