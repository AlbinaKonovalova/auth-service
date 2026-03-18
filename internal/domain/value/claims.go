package value

import (
	"context"
	"strings"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/google/uuid"
)

type ClaimsContextKey struct{}

func NewAccessClaims(
	userID uuid.UUID,
	email string,
	roles []string,
	permissions []string,
) (AccessClaims, error) {
	if userID == uuid.Nil {
		return AccessClaims{}, domain.ErrInvalidUserID
	}
	if strings.TrimSpace(email) == "" {
		return AccessClaims{}, domain.ErrInvalidCredentials
	}

	return AccessClaims{
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

type AccessClaims struct {
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
}

func ClaimsFromContext(ctx context.Context) (AccessClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey{}).(AccessClaims)
	return claims, ok
}
