package domain

import "errors"

var (
	// Auth
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is inactive")

	// User
	ErrUserNotFound = errors.New("user not found")

	// Refresh session
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")

	// Permission
	ErrPermissionNotFound = errors.New("permission not found")
)
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
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters")

	// Role
	ErrRoleNotFound      = errors.New("role not found")
	ErrInvalidRoleCode   = errors.New("invalid role code")
	ErrDuplicateRoleCode = errors.New("duplicate role code")

	// UserRole
	ErrInvalidUserRole = errors.New("invalid user role")

	// Refresh session
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")

	// Permission
	ErrPermissionNotFound = errors.New("permission not found")
)
