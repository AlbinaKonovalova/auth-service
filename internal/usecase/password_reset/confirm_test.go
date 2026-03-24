package password_reset_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	passwordreset "github.com/AlbinaKonovalova/auth-service/internal/usecase/password_reset"
)

// ─── Fixtures ────────────────────────────────────────────────────────────────

var (
	fixedNow    = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	fixedUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	fixedToken  = entity.PasswordResetToken{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserID:    fixedUserID,
		TokenHash: "hashed-token",
		ExpiresAt: fixedNow.Add(30 * time.Minute),
		UsedAt:    nil,
		CreatedAt: fixedNow.Add(-1 * time.Minute),
	}
)

// ─── Mocks ───────────────────────────────────────────────────────────────────

// mockPasswordResetRepo реализует output.PasswordResetRepository
type mockPasswordResetRepo struct {
	findForUpdateFn func(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	findFn          func(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	createFn        func(ctx context.Context, token entity.PasswordResetToken) error
	invalidateFn    func(ctx context.Context, userID uuid.UUID, now time.Time) error
	markUsedFn      func(ctx context.Context, id uuid.UUID, usedAt time.Time) error
}

func (m *mockPasswordResetRepo) FindByTokenHashForUpdate(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	return m.findForUpdateFn(ctx, tokenHash)
}
func (m *mockPasswordResetRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	if m.findFn == nil {
		return nil, domain.ErrResetTokenNotFound
	}
	return m.findFn(ctx, tokenHash)
}
func (m *mockPasswordResetRepo) Create(ctx context.Context, token entity.PasswordResetToken) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(ctx, token)
}
func (m *mockPasswordResetRepo) InvalidateByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error {
	if m.invalidateFn == nil {
		return nil
	}
	return m.invalidateFn(ctx, userID, now)
}
func (m *mockPasswordResetRepo) MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	if m.markUsedFn == nil {
		return nil
	}
	return m.markUsedFn(ctx, id, usedAt)
}

type mockUserRepo struct {
	updatePasswordFn    func(ctx context.Context, id uuid.UUID, hash string) error
	findByIDFn          func(ctx context.Context, id uuid.UUID) (*entity.User, error)
	findByIDForUpdateFn func(ctx context.Context, id uuid.UUID) (*entity.User, error)
	findByEmailFn       func(ctx context.Context, email value.Email) (*entity.User, error)
}

func (m *mockUserRepo) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	if m.updatePasswordFn == nil {
		return nil
	}
	return m.updatePasswordFn(ctx, id, hash)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDFn == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.findByIDFn(ctx, id)
}
func (m *mockUserRepo) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDForUpdateFn == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.findByIDForUpdateFn(ctx, id)
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	if m.findByEmailFn == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.findByEmailFn(ctx, email)
}
func (m *mockUserRepo) Create(_ context.Context, _ entity.User) error   { return nil }
func (m *mockUserRepo) Activate(_ context.Context, _ uuid.UUID) error   { return nil }
func (m *mockUserRepo) Deactivate(_ context.Context, _ uuid.UUID) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _ dto.UserListFilters) ([]entity.User, int, error) {
	return nil, 0, nil
}

type mockSessionRepo struct {
	revokeAllFn func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockSessionRepo) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	if m.revokeAllFn == nil {
		return nil
	}
	return m.revokeAllFn(ctx, userID)
}
func (m *mockSessionRepo) Save(_ context.Context, _ entity.RefreshSession) error { return nil }
func (m *mockSessionRepo) FindByTokenHash(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
	return nil, domain.ErrRefreshTokenNotFound
}
func (m *mockSessionRepo) FindByTokenHashForUpdate(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
	return nil, domain.ErrRefreshTokenNotFound
}
func (m *mockSessionRepo) Revoke(_ context.Context, _ uuid.UUID) error                  { return nil }
func (m *mockSessionRepo) DeleteExpiredAndRevoked(_ context.Context, _ uuid.UUID) error { return nil }

type mockPasswordHasher struct {
	hashFn func(password string) (string, error)
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	if m.hashFn == nil {
		return "new-hashed-" + password, nil
	}
	return m.hashFn(password)
}
func (m *mockPasswordHasher) Verify(_, _ string) (bool, error) { return true, nil }

type mockTokenProvider struct{}

func (m *mockTokenProvider) GenerateAccessToken(_ context.Context, _ value.AccessClaims) (string, int64, error) {
	return "", 0, nil
}
func (m *mockTokenProvider) ParseAccessToken(_ context.Context, _ string) (value.AccessClaims, error) {
	return value.AccessClaims{}, nil
}
func (m *mockTokenProvider) GenerateRefreshToken(_ context.Context) (string, value.TokenHash, error) {
	return "", "", nil
}
func (m *mockTokenProvider) GeneratePasswordResetToken(_ context.Context) (string, value.TokenHash, error) {
	return "raw-token", "hashed-token", nil
}

type mockTokenHasher struct{}

func (m *mockTokenHasher) Hash(raw string) value.TokenHash {
	// детерминированный хеш для тестов: просто префикс
	return value.TokenHash("hashed-" + raw)
}

type mockMailer struct {
	sendFn func(ctx context.Context, toEmail, resetToken string) error
}

func (m *mockMailer) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	if m.sendFn == nil {
		return nil
	}
	return m.sendFn(ctx, toEmail, resetToken)
}

type mockClock struct{ now time.Time }

func (m *mockClock) Now() time.Time { return m.now }

// mockUUIDGen реализует output.UUIDGenerator
type mockUUIDGen struct{ id uuid.UUID }

func (m *mockUUIDGen) New() uuid.UUID { return m.id }

type inlineTxManager struct {
	// rollback: если fn вернула ошибку, txManager её пробрасывает как есть (как реальный Postgres).
}

func (t *inlineTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ─── Builder ─────────────────────────────────────────────────────────────────

type confirmBuilder struct {
	resetRepo   *mockPasswordResetRepo
	userRepo    *mockUserRepo
	sessionRepo *mockSessionRepo
	hasher      *mockPasswordHasher
	clock       *mockClock
}

func newConfirmBuilder() *confirmBuilder {
	return &confirmBuilder{
		resetRepo: &mockPasswordResetRepo{
			findForUpdateFn: func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
				tok := fixedToken // копия
				return &tok, nil
			},
			markUsedFn: func(_ context.Context, _ uuid.UUID, _ time.Time) error {
				return nil
			},
		},
		userRepo: &mockUserRepo{
			updatePasswordFn: func(_ context.Context, _ uuid.UUID, _ string) error {
				return nil
			},
		},
		sessionRepo: &mockSessionRepo{
			revokeAllFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		},
		hasher: &mockPasswordHasher{},
		clock:  &mockClock{now: fixedNow},
	}
}

func (b *confirmBuilder) build() *passwordreset.PasswordResetService {
	return passwordreset.NewPasswordResetService(
		b.userRepo,
		b.resetRepo,
		b.sessionRepo,
		b.hasher,
		&mockTokenProvider{},
		&mockTokenHasher{},
		&mockMailer{},
		b.clock,
		&mockUUIDGen{id: uuid.New()},
		&inlineTxManager{},
		30*time.Minute,
	)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestConfirmPasswordReset_Success(t *testing.T) {
	var (
		passwordUpdated bool
		tokenMarkedUsed bool
		sessionsRevoked bool
	)

	b := newConfirmBuilder()

	b.resetRepo.findForUpdateFn = func(_ context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
		if tokenHash != "hashed-raw-token" {
			t.Errorf("unexpected token hash: got %q, want %q", tokenHash, "hashed-raw-token")
		}
		tok := fixedToken
		return &tok, nil
	}
	b.userRepo.updatePasswordFn = func(_ context.Context, id uuid.UUID, hash string) error {
		if id != fixedUserID {
			t.Errorf("unexpected user id: %v", id)
		}
		if hash == "" {
			t.Error("password hash must not be empty")
		}
		passwordUpdated = true
		return nil
	}
	b.resetRepo.markUsedFn = func(_ context.Context, id uuid.UUID, _ time.Time) error {
		if id != fixedToken.ID {
			t.Errorf("unexpected token id: %v", id)
		}
		tokenMarkedUsed = true
		return nil
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, userID uuid.UUID) error {
		if userID != fixedUserID {
			t.Errorf("unexpected user id for revoke: %v", userID)
		}
		sessionsRevoked = true
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !passwordUpdated {
		t.Error("expected password hash to be updated")
	}
	if !tokenMarkedUsed {
		t.Error("expected token to be marked as used")
	}
	if !sessionsRevoked {
		t.Error("expected refresh sessions to be revoked")
	}
}

func TestConfirmPasswordReset_ExpiredToken(t *testing.T) {
	expiredNow := fixedToken.ExpiresAt.Add(1 * time.Second) // после ExpiresAt

	b := newConfirmBuilder()
	b.clock = &mockClock{now: expiredNow}

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		tok := fixedToken
		return &tok, nil
	}
	b.userRepo.updatePasswordFn = func(_ context.Context, _ uuid.UUID, _ string) error {
		t.Error("UpdatePasswordHash must not be called for expired token")
		return nil
	}
	b.resetRepo.markUsedFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		t.Error("MarkUsed must not be called for expired token")
		return nil
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("RevokeAllByUserID must not be called for expired token")
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if !errors.Is(err, domain.ErrResetTokenExpired) {
		t.Fatalf("expected ErrResetTokenExpired, got: %v", err)
	}
}

func TestConfirmPasswordReset_UsedToken(t *testing.T) {
	usedAt := fixedNow.Add(-5 * time.Minute)

	b := newConfirmBuilder()

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		tok := fixedToken
		tok.UsedAt = &usedAt // token уже использован
		return &tok, nil
	}
	b.userRepo.updatePasswordFn = func(_ context.Context, _ uuid.UUID, _ string) error {
		t.Error("UpdatePasswordHash must not be called for used token")
		return nil
	}
	b.resetRepo.markUsedFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		t.Error("MarkUsed must not be called for used token")
		return nil
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("RevokeAllByUserID must not be called for used token")
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if !errors.Is(err, domain.ErrResetTokenUsed) {
		t.Fatalf("expected ErrResetTokenUsed, got: %v", err)
	}
}

func TestConfirmPasswordReset_UnknownToken(t *testing.T) {
	b := newConfirmBuilder()

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		return nil, domain.ErrResetTokenNotFound
	}
	b.userRepo.updatePasswordFn = func(_ context.Context, _ uuid.UUID, _ string) error {
		t.Error("UpdatePasswordHash must not be called for unknown token")
		return nil
	}
	b.resetRepo.markUsedFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		t.Error("MarkUsed must not be called for unknown token")
		return nil
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("RevokeAllByUserID must not be called for unknown token")
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "nonexistent-token",
		NewPassword: "ValidPassword1!",
	})

	if !errors.Is(err, domain.ErrResetTokenNotFound) {
		t.Fatalf("expected ErrResetTokenNotFound, got: %v", err)
	}
}

func TestConfirmPasswordReset_InvalidPassword(t *testing.T) {
	b := newConfirmBuilder()

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		t.Error("FindByTokenHashForUpdate must not be called when password is invalid")
		return nil, nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "short", // менее 8 символов
	})

	if !errors.Is(err, domain.ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got: %v", err)
	}
}

func TestConfirmPasswordReset_RevokeSessionsAfterSuccess(t *testing.T) {
	var callOrder []string

	b := newConfirmBuilder()

	b.userRepo.updatePasswordFn = func(_ context.Context, _ uuid.UUID, _ string) error {
		callOrder = append(callOrder, "update_password")
		return nil
	}
	b.resetRepo.markUsedFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		callOrder = append(callOrder, "mark_used")
		return nil
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		callOrder = append(callOrder, "revoke_sessions")
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedOrder := []string{"update_password", "mark_used", "revoke_sessions"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("expected %d calls, got %d: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, step := range expectedOrder {
		if callOrder[i] != step {
			t.Errorf("step %d: expected %q, got %q", i, step, callOrder[i])
		}
	}
}

func TestConfirmPasswordReset_RollbackOnMarkUsedFailure(t *testing.T) {
	var sessionsRevoked bool

	b := newConfirmBuilder()

	b.resetRepo.markUsedFn = func(_ context.Context, _ uuid.UUID, _ time.Time) error {
		// симулирует: второй параллельный confirm — token уже помечен первым
		return domain.ErrResetTokenUsed
	}
	b.sessionRepo.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		sessionsRevoked = true
		return nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if !errors.Is(err, domain.ErrResetTokenUsed) {
		t.Fatalf("expected ErrResetTokenUsed, got: %v", err)
	}
	if sessionsRevoked {
		t.Error("sessions must not be revoked if MarkUsed failed")
	}
}

func TestConfirmPasswordReset_ParallelConflictCaughtByLock(t *testing.T) {
	usedAt := fixedNow

	callCount := 0

	b := newConfirmBuilder()

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		callCount++
		tok := fixedToken
		if callCount > 1 {
			tok.UsedAt = &usedAt
		}
		return &tok, nil
	}

	svc := b.build()

	err1 := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})
	if err1 != nil {
		t.Fatalf("first confirm: expected no error, got: %v", err1)
	}

	err2 := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})
	if !errors.Is(err2, domain.ErrResetTokenUsed) {
		t.Fatalf("second confirm: expected ErrResetTokenUsed, got: %v", err2)
	}
}

func TestConfirmPasswordReset_ExpiredCheckedBeforeUsed(t *testing.T) {
	usedAt := fixedNow.Add(-10 * time.Minute)
	expiredNow := fixedToken.ExpiresAt.Add(1 * time.Second)

	b := newConfirmBuilder()
	b.clock = &mockClock{now: expiredNow}

	b.resetRepo.findForUpdateFn = func(_ context.Context, _ string) (*entity.PasswordResetToken, error) {
		tok := fixedToken
		tok.UsedAt = &usedAt // и использован, и истёк
		return &tok, nil
	}

	svc := b.build()
	err := svc.ConfirmPasswordReset(context.Background(), input.ConfirmPasswordResetInput{
		Token:       "raw-token",
		NewPassword: "ValidPassword1!",
	})

	if !errors.Is(err, domain.ErrResetTokenExpired) {
		t.Fatalf("expected ErrResetTokenExpired (checked before used), got: %v", err)
	}
}
