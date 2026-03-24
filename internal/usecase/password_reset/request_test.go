package password_reset_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	passwordreset "github.com/AlbinaKonovalova/auth-service/internal/usecase/password_reset"
)

// ─── TestRequestPasswordReset ─────────────────────────────────────────────────

func TestRequestPasswordReset_Success(t *testing.T) {
	var (
		oldInvalidated bool
		tokenCreated   bool
		emailSent      bool
	)

	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		u := entity.User{ID: fixedUserID, Email: "admin@example.com", IsActive: true}
		return &u, nil
	}
	b.resetRepo.invalidateFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		oldInvalidated = true
		return nil
	}
	b.resetRepo.createFn = func(_ context.Context, _ entity.PasswordResetToken) error {
		tokenCreated = true
		return nil
	}
	b.mailer.sendFn = func(_ context.Context, _, _ string) error {
		emailSent = true
		return nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "admin@example.com",
	})

	require.NoError(t, err)
	assert.True(t, oldInvalidated, "old tokens must be invalidated")
	assert.True(t, tokenCreated, "new token must be created")
	assert.True(t, emailSent, "email must be sent")
}

func TestRequestPasswordReset_UnknownEmail_SilentSuccess(t *testing.T) {
	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		return nil, domain.ErrUserNotFound
	}
	b.resetRepo.createFn = func(_ context.Context, _ entity.PasswordResetToken) error {
		t.Error("must not create token for unknown email")
		return nil
	}
	b.mailer.sendFn = func(_ context.Context, _, _ string) error {
		t.Error("must not send email for unknown email")
		return nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "nobody@example.com",
	})
	assert.NoError(t, err) // no user enumeration
}

func TestRequestPasswordReset_InvalidEmail_SilentSuccess(t *testing.T) {
	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		t.Error("must not query DB for invalid email")
		return nil, nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "not-an-email",
	})
	assert.NoError(t, err)
}

func TestRequestPasswordReset_DBError_Propagated(t *testing.T) {
	infraErr := errors.New("connection lost")

	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		return nil, infraErr
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "admin@example.com",
	})
	assert.ErrorContains(t, err, "connection lost")
}

func TestRequestPasswordReset_OldTokensInvalidatedBeforeNew(t *testing.T) {
	var callOrder []string

	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		u := entity.User{ID: fixedUserID, Email: "admin@example.com", IsActive: true}
		return &u, nil
	}
	b.resetRepo.invalidateFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		callOrder = append(callOrder, "invalidate")
		return nil
	}
	b.resetRepo.createFn = func(_ context.Context, _ entity.PasswordResetToken) error {
		callOrder = append(callOrder, "create")
		return nil
	}
	b.mailer.sendFn = func(_ context.Context, _, _ string) error {
		callOrder = append(callOrder, "send_email")
		return nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "admin@example.com",
	})

	require.NoError(t, err)
	require.Equal(t, []string{"invalidate", "create", "send_email"}, callOrder)
}

func TestRequestPasswordReset_EmailNotSentIfTxFails(t *testing.T) {
	var emailSent bool
	txErr := errors.New("tx failed")

	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		u := entity.User{ID: fixedUserID, Email: "admin@example.com", IsActive: true}
		return &u, nil
	}
	b.resetRepo.createFn = func(_ context.Context, _ entity.PasswordResetToken) error {
		return txErr // транзакция упала
	}
	b.mailer.sendFn = func(_ context.Context, _, _ string) error {
		emailSent = true
		return nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "admin@example.com",
	})
	assert.Error(t, err)
	assert.False(t, emailSent, "email must not be sent if transaction failed")
}

func TestRequestPasswordReset_TokenDataContract(t *testing.T) {
	var createdTokenHash string
	var mailedRawToken string

	b := newRequestBuilder()
	b.userRepo.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		u := entity.User{ID: fixedUserID, Email: "admin@example.com", IsActive: true}
		return &u, nil
	}
	b.resetRepo.createFn = func(_ context.Context, tok entity.PasswordResetToken) error {
		createdTokenHash = tok.TokenHash
		return nil
	}
	b.mailer.sendFn = func(_ context.Context, _, resetToken string) error {
		mailedRawToken = resetToken
		return nil
	}

	svc := b.build()
	err := svc.RequestPasswordReset(context.Background(), input.RequestPasswordResetInput{
		Email: "admin@example.com",
	})
	require.NoError(t, err)

	assert.Equal(t, "hashed-token", createdTokenHash, "repo должен получить hash, а не raw token")
	assert.Equal(t, "raw-token", mailedRawToken, "mailer должен получить raw token, а не hash")
	assert.NotEqual(t, createdTokenHash, mailedRawToken, "hash и raw token не должны совпадать")
}

type requestBuilder struct {
	userRepo  *mockUserRepo
	resetRepo *mockPasswordResetRepo
	mailer    *mockMailer
}

func newRequestBuilder() *requestBuilder {
	b := &requestBuilder{
		userRepo:  &mockUserRepo{},
		resetRepo: &mockPasswordResetRepo{},
		mailer:    &mockMailer{},
	}

	b.userRepo.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		u := entity.User{ID: fixedUserID, Email: "admin@example.com", IsActive: true}
		return &u, nil
	}
	return b
}

func (b *requestBuilder) build() *passwordreset.PasswordResetService {
	return passwordreset.NewPasswordResetService(
		b.userRepo,
		b.resetRepo,
		&mockSessionRepo{},
		&mockPasswordHasher{},
		&mockTokenProvider{},
		&mockTokenHasher{},
		b.mailer,
		&mockClock{now: fixedNow},
		&mockUUIDGen{id: uuid.New()},
		&inlineTxManager{},
		30*time.Minute,
	)
}
