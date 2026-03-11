package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/AlbinaKonovalova/auth-service/internal/config"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

var errCookieMissing = errors.New("cookie missing")

var errCookieMalformed = errors.New("cookie malformed")

type AuthServiceController struct {
	pb.UnimplementedAuthServiceServer
	auth   input.AuthUseCase
	cookie config.CookieConfig
}

func NewAuthServiceController(
	auth input.AuthUseCase,
	cookie config.CookieConfig,
) *AuthServiceController {
	return &AuthServiceController{
		auth:   auth,
		cookie: cookie,
	}
}

func (c *AuthServiceController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, rawRefresh, err := c.auth.Login(ctx, input.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, domainErrToStatus(err)
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
		return nil, domainErrToStatus(err)
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
			return nil, domainErrToStatus(logoutErr)
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
		return nil, domainErrToStatus(err)
	}

	return &pb.MeResponse{
		Id:          user.ID.String(),
		Email:       user.Email,
		IsActive:    user.IsActive,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}, nil
}

func (c *AuthServiceController) Healthz(_ context.Context, _ *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (c *AuthServiceController) setRefreshCookie(ctx context.Context, raw string) error {
	cookie := &http.Cookie{
		Name:     c.cookie.Name,
		Value:    raw,
		Path:     c.cookie.Path,
		Domain:   c.cookie.Domain,
		HttpOnly: true,
		Secure:   c.cookie.IsSecure(),
		SameSite: parseSameSite(c.cookie.SameSite),
		MaxAge:   int(c.cookie.TTL.Seconds()),
	}
	return grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
}

func (c *AuthServiceController) clearRefreshCookie(ctx context.Context) error {
	cookie := &http.Cookie{
		Name:     c.cookie.Name,
		Value:    "",
		Path:     c.cookie.Path,
		Domain:   c.cookie.Domain,
		HttpOnly: true,
		Secure:   c.cookie.IsSecure(),
		SameSite: parseSameSite(c.cookie.SameSite),
		MaxAge:   -1,
	}
	return grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
}

func (c *AuthServiceController) readRefreshCookie(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("%w: no incoming metadata", errCookieMalformed)
	}

	cookies := md.Get("cookie")
	if len(cookies) == 0 {
		return "", errCookieMissing
	}

	header := http.Header{}
	for _, v := range cookies {
		header.Add("Cookie", v)
	}
	r := &http.Request{Header: header}

	cookie, err := r.Cookie(c.cookie.Name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errCookieMissing
		}
		return "", fmt.Errorf("%w: %w", errCookieMalformed, err)
	}

	return cookie.Value, nil
}

func parseSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteStrictMode
	}
}
