package auth

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
)

// AuthService объединяет отдельные usecase в единую реализацию input.AuthUseCase.
type AuthService struct {
	login   *LoginUseCase
	refresh *RefreshUseCase
	logout  *LogoutUseCase
	me      *MeUseCase
}

func NewAuthService(
	login *LoginUseCase,
	refresh *RefreshUseCase,
	logout *LogoutUseCase,
	me *MeUseCase,
) *AuthService {
	return &AuthService{
		login:   login,
		refresh: refresh,
		logout:  logout,
		me:      me,
	}
}

func (s *AuthService) Login(ctx context.Context, in input.LoginInput) (dto.LoginResult, string, error) {
	return s.login.Login(ctx, in)
}

func (s *AuthService) Refresh(ctx context.Context, in input.RefreshInput) (dto.RefreshResult, string, error) {
	return s.refresh.Refresh(ctx, in)
}

func (s *AuthService) Logout(ctx context.Context, in input.LogoutInput) error {
	return s.logout.Logout(ctx, in)
}

func (s *AuthService) Me(ctx context.Context, claims value.AccessClaims) (dto.CurrentUser, error) {
	return s.me.Me(ctx, claims)
}

// compile-time check
var _ input.AuthUseCase = (*AuthService)(nil)
