package controller

import (
	"errors"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AuthServiceController) domainErrToStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, domain.ErrUserInactive):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, domain.ErrRoleInUse):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrEmailAlreadyTaken),
		errors.Is(err, domain.ErrDuplicateRoleCode),
		errors.Is(err, domain.ErrDuplicatePermissionCode):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, domain.ErrRoleNotFound),
		errors.Is(err, domain.ErrUserRoleNotFound),
		errors.Is(err, domain.ErrRolePermissionNotFound),
		errors.Is(err, domain.ErrPermissionNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrInvalidRoleCode),
		errors.Is(err, domain.ErrRoleNameEmpty),
		errors.Is(err, domain.ErrInvalidUserID),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidPasswordHash),
		errors.Is(err, domain.ErrInvalidPassword),
		errors.Is(err, domain.ErrInvalidUserRole),
		errors.Is(err, domain.ErrInvalidRolePermission),
		errors.Is(err, domain.ErrInvalidUserList),
		errors.Is(err, domain.ErrInvalidPermissionCode),
		errors.Is(err, domain.ErrPermissionDescriptionEmpty),
		errors.Is(err, domain.ErrUserMustHaveRole),
		errors.Is(err, domain.ErrCannotRevokeLastRole):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, domain.ErrRefreshTokenNotFound),
		errors.Is(err, domain.ErrRefreshTokenRevoked),
		errors.Is(err, domain.ErrRefreshTokenExpired):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, domain.ErrDataIntegrityViolation):
		c.logger.Error("data integrity violation", "error", err)
		return status.Error(codes.Internal, "internal server error")

	default:
		c.logger.Error("unexpected error", "error", err)
		return status.Error(codes.Internal, "internal server error")
	}
}
