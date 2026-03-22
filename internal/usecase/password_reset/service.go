package password_reset

import (
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

// PasswordResetService реализует input.PasswordResetUseCase.
// Оркестрирует сценарии request и confirm password reset.
type PasswordResetService struct {
	userRepo      output.UserRepository
	resetRepo     output.PasswordResetRepository
	sessionRepo   output.RefreshSessionRepository
	hasher        output.PasswordHasher
	tokenProvider output.TokenProvider
	tokenHasher   output.TokenHasher
	mailer        output.Mailer
	clock         output.Clock
	uuidGen       output.UUIDGenerator
	tx            output.TxManager
	tokenTTL      time.Duration
}

func NewPasswordResetService(
	userRepo output.UserRepository,
	resetRepo output.PasswordResetRepository,
	sessionRepo output.RefreshSessionRepository,
	hasher output.PasswordHasher,
	tokenProvider output.TokenProvider,
	tokenHasher output.TokenHasher,
	mailer output.Mailer,
	clock output.Clock,
	uuidGen output.UUIDGenerator,
	tx output.TxManager,
	tokenTTL time.Duration,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:      userRepo,
		resetRepo:     resetRepo,
		sessionRepo:   sessionRepo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		tokenHasher:   tokenHasher,
		mailer:        mailer,
		clock:         clock,
		uuidGen:       uuidGen,
		tx:            tx,
		tokenTTL:      tokenTTL,
	}
}
