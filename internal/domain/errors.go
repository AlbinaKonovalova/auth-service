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
