package service

import (
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

// ValidateResetToken проверяет, что токен применим в момент now.
// Это единственный источник истины для бизнес-правила usability reset token.
// Делегирует в entity.PasswordResetToken.EnsureUsable — доменное правило живёт на сущности,
// а этот helper позволяет вызывать его единообразно из usecase без прямой зависимости от entity-метода.
//
// Возвращает domain.ErrResetTokenExpired если токен истёк.
// Возвращает domain.ErrResetTokenUsed если токен уже был использован или аннулирован.
func ValidateResetToken(token *entity.PasswordResetToken, now time.Time) error {
	return token.EnsureUsable(now)
}
