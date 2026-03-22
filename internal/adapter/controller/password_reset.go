package controller

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestPasswordReset инициирует password reset flow.
// Всегда возвращает success независимо от того, существует email или нет.
// Ошибки инфраструктуры (DB, SMTP) возвращаются как internal.
func (c *AuthServiceController) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	err := c.passwordReset.RequestPasswordReset(ctx, input.RequestPasswordResetInput{
		Email: req.Email,
	})
	if err != nil {
		c.logger.Error("password reset request failed", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.RequestPasswordResetResponse{Success: true}, nil
}

// ConfirmPasswordReset завершает password reset flow.
// Доменные ошибки (not found, expired, used, invalid password) маппятся в transport статусы.
func (c *AuthServiceController) ConfirmPasswordReset(ctx context.Context, req *pb.ConfirmPasswordResetRequest) (*pb.ConfirmPasswordResetResponse, error) {
	err := c.passwordReset.ConfirmPasswordReset(ctx, input.ConfirmPasswordResetInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.ConfirmPasswordResetResponse{Success: true}, nil
}
