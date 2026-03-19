package domain

import "errors"

var (
	// Auth
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is inactive")

	// User
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyTaken   = errors.New("email already taken")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrUserMustHaveRole    = errors.New("user must have at least one role")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters")
	ErrInvalidUserList     = errors.New("invalid user list")

	// Role
	ErrRoleNotFound      = errors.New("role not found")
	ErrInvalidRoleCode   = errors.New("invalid role code")
	ErrDuplicateRoleCode = errors.New("duplicate role code")
	ErrRoleNameEmpty     = errors.New("role name must not be empty")

	// UserRole
	ErrInvalidUserRole = errors.New("invalid user role")

	// RolePermission
	ErrInvalidRolePermission = errors.New("invalid role permission")

	// Refresh session
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")

	// Permission
	ErrPermissionNotFound    = errors.New("permission not found")
	ErrInvalidPermissionCode = errors.New("invalid permission code")

	// Access — assign / revoke
	ErrCannotRevokeLastRole   = errors.New("cannot revoke the last role from user")
	ErrUserRoleNotFound       = errors.New("user does not have this role")
	ErrRolePermissionNotFound = errors.New("role does not have this permission")

	// Data integrity — битые связи в БД, не бизнес-ошибки запроса пользователя
	ErrDataIntegrityViolation = errors.New("data integrity violation")
)
