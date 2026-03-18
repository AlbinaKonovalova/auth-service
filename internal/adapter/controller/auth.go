package controller

import (
	"context"
	"errors"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AuthServiceController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, rawRefresh, err := c.auth.Login(ctx, input.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	if err := c.setRefreshCookie(ctx, rawRefresh); err != nil {
		return nil, status.Error(codes.Internal, "failed to set cookie")
	}

	return loginResultToProto(result), nil
}

func (c *AuthServiceController) Refresh(ctx context.Context, _ *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	rawRefresh, err := c.readRefreshCookie(ctx)
	if err != nil {
		if errors.Is(err, errCookieMalformed) {
			return nil, status.Error(codes.InvalidArgument, "malformed cookie header")
		}
		return nil, status.Error(codes.Unauthenticated, "refresh token cookie missing")
	}

	result, newRawRefresh, err := c.auth.Refresh(ctx, input.RefreshInput{
		RawRefreshToken: rawRefresh,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	if err := c.setRefreshCookie(ctx, newRawRefresh); err != nil {
		return nil, status.Error(codes.Internal, "failed to set cookie")
	}

	return refreshResultToProto(result), nil
}

func (c *AuthServiceController) Logout(ctx context.Context, _ *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	rawRefresh, err := c.readRefreshCookie(ctx)
	switch {
	case err == nil:
		if logoutErr := c.auth.Logout(ctx, input.LogoutInput{
			RawRefreshToken: rawRefresh,
		}); logoutErr != nil {
			return nil, c.domainErrToStatus(logoutErr)
		}
	case errors.Is(err, errCookieMissing):
	default:
		return nil, status.Error(codes.InvalidArgument, "malformed cookie header")
	}

	if err := c.clearRefreshCookie(ctx); err != nil {
		return nil, status.Error(codes.Internal, "failed to clear cookie")
	}

	return &pb.LogoutResponse{Success: true}, nil
}

func (c *AuthServiceController) Me(ctx context.Context, _ *pb.MeRequest) (*pb.MeResponse, error) {
	claims, ok := ctx.Value(value.ClaimsContextKey{}).(value.AccessClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth claims")
	}

	user, err := c.auth.Me(ctx, claims)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.MeResponse{
		Id:          user.ID.String(),
		Email:       user.Email,
		IsActive:    user.IsActive,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}, nil
}
