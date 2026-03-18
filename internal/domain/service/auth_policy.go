package service

import (
	"errors"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

func MapLoginLookupError(err error) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return domain.ErrInvalidCredentials
	}

	return err
}

func CanLogout(err error, session *entity.RefreshSession) (bool, error) {
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return false, nil
		}
		return false, err
	}

	if session == nil {
		return false, nil
	}

	if session.IsRevoked() {
		return false, nil
	}

	return true, nil
}
