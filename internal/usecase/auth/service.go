package auth

import (
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

const refreshTTL = 7 * 24 * time.Hour

type AuthService struct {
	users    output.UserRepository
	sessions output.RefreshSessionRepository
	hasher   output.PasswordHasher
	tokens   output.TokenProvider
	tokenH   output.TokenHasher
	tx       output.TxManager
	clock    output.Clock
	uuid     output.UUIDGenerator
	resolver *common.PermissionResolver
}

func NewAuthService(
	users output.UserRepository,
	sessions output.RefreshSessionRepository,
	hasher output.PasswordHasher,
	tokens output.TokenProvider,
	tokenH output.TokenHasher,
	tx output.TxManager,
	clock output.Clock,
	uuid output.UUIDGenerator,
	resolver *common.PermissionResolver,
) *AuthService {
	return &AuthService{
		users:    users,
		sessions: sessions,
		hasher:   hasher,
		tokens:   tokens,
		tokenH:   tokenH,
		tx:       tx,
		clock:    clock,
		uuid:     uuid,
		resolver: resolver,
	}
}

// compile-time check
var _ input.AuthUseCase = (*AuthService)(nil)
