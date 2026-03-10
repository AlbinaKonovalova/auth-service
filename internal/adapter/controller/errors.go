package controller

import (
	"errors"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// domainErrToStatus maps domain errors to gRPC status codes.
// Example:
//
//	domain.ErrInvalidCredentials -> codes.Unauthenticated
//	domain.ErrNotFound           -> codes.NotFound
//	domain.ErrAccessDenied       -> codes.PermissionDenied
func domainErrToStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrUserInactive):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrRefreshTokenNotFound),
		errors.Is(err, domain.ErrRefreshTokenRevoked),
		errors.Is(err, domain.ErrRefreshTokenExpired):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrPermissionNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
