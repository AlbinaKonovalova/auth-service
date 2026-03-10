package value

import "github.com/google/uuid"

// AccessClaims — доменная модель claims access token.
// Не привязана к jwt library.
type AccessClaims struct {
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
}

// ClaimsContextKey — ключ для хранения AccessClaims в context.Context.
type ClaimsContextKey struct{}
