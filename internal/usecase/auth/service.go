package auth

import (
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

type AuthService struct {
	users      output.UserRepository
	sessions   output.RefreshSessionRepository
	hasher     output.PasswordHasher
	tokens     output.TokenProvider
	tokenH     output.TokenHasher
	tx         output.TxManager
	clock      output.Clock
	uuid       output.UUIDGenerator
	resolver   *common.PermissionResolver
	refreshTTL time.Duration
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
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users:      users,
		sessions:   sessions,
		hasher:     hasher,
		tokens:     tokens,
		tokenH:     tokenH,
		tx:         tx,
		clock:      clock,
		uuid:       uuid,
		resolver:   resolver,
		refreshTTL: refreshTTL,
	}
}

var _ input.AuthUseCase = (*AuthService)(nil)
