package value

import "github.com/google/uuid"

type AccessClaims struct {
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
}

type ClaimsContextKey struct{}
