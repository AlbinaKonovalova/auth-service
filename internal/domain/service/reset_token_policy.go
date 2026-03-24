package service

import (
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func ValidateResetToken(token *entity.PasswordResetToken, now time.Time) error {
	return token.EnsureUsable(now)
}
