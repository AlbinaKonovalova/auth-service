package input

import "context"

type RequestPasswordResetInput struct {
	Email string
}

type ConfirmPasswordResetInput struct {
	// Token — сырой токен из ссылки в письме (не хеш).
	Token       string
	NewPassword string
}

type PasswordResetUseCase interface {
	// RequestPasswordReset инициирует сброс пароля.
	// На уровне foundation (commit 13) это контракт сценария.
	// Реализация (создание токена, отправка письма) идёт в commit 14.
	RequestPasswordReset(ctx context.Context, in RequestPasswordResetInput) error

	ConfirmPasswordReset(ctx context.Context, in ConfirmPasswordResetInput) error
}
