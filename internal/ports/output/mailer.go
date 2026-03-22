package output

import "context"

// Mailer — контракт для отправки email-сообщений.
type Mailer interface {
	// SendPasswordResetEmail отправляет письмо со ссылкой для сброса пароля.
	// toEmail — адрес получателя.
	// resetToken — сырой (не захешированный) токен, который вставляется в ссылку.
	SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error
}
