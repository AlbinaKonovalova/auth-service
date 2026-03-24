package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is inactive")

	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyTaken   = errors.New("email already taken")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrUserMustHaveRole    = errors.New("user must have at least one role")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters")
	ErrInvalidUserList     = errors.New("invalid user list")

	ErrRoleNotFound      = errors.New("role not found")
	ErrInvalidRoleCode   = errors.New("invalid role code")
	ErrDuplicateRoleCode = errors.New("duplicate role code")
	ErrRoleNameEmpty     = errors.New("role name must not be empty")
	ErrRoleInUse         = errors.New("role is assigned and cannot be deleted")

	ErrInvalidUserRole = errors.New("invalid user role")

	ErrInvalidRolePermission = errors.New("invalid role permission")

	ErrRefreshTokenNotFound      = errors.New("refresh token not found")
	ErrRefreshTokenRevoked       = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired       = errors.New("refresh token has expired")
	ErrInvalidRefreshSessionID   = errors.New("invalid refresh session id")
	ErrInvalidRefreshTokenHash   = errors.New("invalid refresh token hash")
	ErrInvalidRefreshTokenTTL    = errors.New("invalid refresh token ttl")
	ErrInvalidRefreshSessionTime = errors.New("invalid refresh session time")

	ErrPermissionNotFound         = errors.New("permission not found")
	ErrInvalidPermissionCode      = errors.New("invalid permission code")
	ErrDuplicatePermissionCode    = errors.New("duplicate permission code")
	ErrPermissionDescriptionEmpty = errors.New("permission description must not be empty")
	ErrPermissionInUse            = errors.New("permission is assigned to a role and cannot be deleted")

	ErrCannotRevokeLastRole   = errors.New("cannot revoke the last role from user")
	ErrUserRoleNotFound       = errors.New("user does not have this role")
	ErrRolePermissionNotFound = errors.New("role does not have this permission")

	ErrDataIntegrityViolation = errors.New("data integrity violation")

	ErrResetTokenNotFound      = errors.New("reset token not found")
	ErrResetTokenExpired       = errors.New("reset token has expired")
	ErrResetTokenUsed          = errors.New("reset token has already been used")
	ErrInvalidResetTokenID     = errors.New("invalid reset token id")
	ErrInvalidResetTokenUserID = errors.New("invalid reset token user id")
	ErrInvalidResetTokenHash   = errors.New("invalid reset token hash")
	ErrInvalidResetTokenNow    = errors.New("invalid reset token time")
	ErrInvalidResetTokenTTL    = errors.New("invalid reset token ttl")
)
