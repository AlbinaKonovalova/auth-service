package controller

import (
	"context"
	"log/slog"

	"github.com/AlbinaKonovalova/auth-service/internal/config/modules"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

type AuthServiceController struct {
	pb.UnimplementedAuthServiceServer
	auth          input.AuthUseCase
	user          input.UserUseCase
	access        input.AccessUseCase
	role          input.RoleUseCase
	permission    input.PermissionUseCase
	passwordReset input.PasswordResetUseCase
	cookie        modules.CookieConfig
	logger        *slog.Logger
}

func NewAuthServiceController(
	auth input.AuthUseCase,
	user input.UserUseCase,
	access input.AccessUseCase,
	role input.RoleUseCase,
	permission input.PermissionUseCase,
	passwordReset input.PasswordResetUseCase,
	cookie modules.CookieConfig,
	logger *slog.Logger,
) *AuthServiceController {
	return &AuthServiceController{
		auth:          auth,
		user:          user,
		access:        access,
		role:          role,
		permission:    permission,
		passwordReset: passwordReset,
		cookie:        cookie,
		logger:        logger,
	}
}

func (c *AuthServiceController) Healthz(_ context.Context, _ *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{Status: "ok"}, nil
}
