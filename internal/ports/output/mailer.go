package output

import "context"

type Mailer interface {
	SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error
}
