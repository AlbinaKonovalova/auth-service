package input

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type LoginInput struct {
	Email    string
	Password string
}

type RefreshInput struct {
	RawRefreshToken string
}

type LogoutInput struct {
	RawRefreshToken string
}

type AuthUseCase interface {
	Login(ctx context.Context, in LoginInput) (result dto.LoginResult, rawRefreshToken string, err error)
	Refresh(ctx context.Context, in RefreshInput) (result dto.RefreshResult, rawRefreshToken string, err error)
	Logout(ctx context.Context, in LogoutInput) error
	Me(ctx context.Context, claims value.AccessClaims) (dto.CurrentUser, error)
}
