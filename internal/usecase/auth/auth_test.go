package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	authusecase "github.com/AlbinaKonovalova/auth-service/internal/usecase/auth"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
)

// ─── Fixtures ─────────────────────────────────────────────────────────────────

var (
	testNow    = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	testUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	activeUser = entity.User{
		ID:           testUserID,
		Email:        "admin@example.com",
		PasswordHash: "hashed_password",
		IsActive:     true,
	}
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	findByEmailFn func(ctx context.Context, email value.Email) (*entity.User, error)
	findByIDFn    func(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	if m.findByEmailFn == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.findByEmailFn(ctx, email)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDFn == nil {
		return nil, domain.ErrUserNotFound
	}
	return m.findByIDFn(ctx, id)
}
func (m *mockUserRepo) FindByIDForUpdate(_ context.Context, _ uuid.UUID) (*entity.User, error) {
	return nil, domain.ErrUserNotFound
}
func (m *mockUserRepo) Create(_ context.Context, _ entity.User) error                     { return nil }
func (m *mockUserRepo) UpdatePasswordHash(_ context.Context, _ uuid.UUID, _ string) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _ dto.UserListFilters) ([]entity.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Activate(_ context.Context, _ uuid.UUID) error   { return nil }
func (m *mockUserRepo) Deactivate(_ context.Context, _ uuid.UUID) error { return nil }

type mockSessionRepo struct {
	saveFn                  func(ctx context.Context, s entity.RefreshSession) error
	findByTokenHashFn       func(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error)
	findByTokenHashForUpdFn func(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error)
	revokeFn                func(ctx context.Context, id uuid.UUID) error
	deleteExpiredFn         func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockSessionRepo) Save(ctx context.Context, s entity.RefreshSession) error {
	if m.saveFn == nil {
		return nil
	}
	return m.saveFn(ctx, s)
}
func (m *mockSessionRepo) FindByTokenHash(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error) {
	if m.findByTokenHashFn == nil {
		return nil, domain.ErrRefreshTokenNotFound
	}
	return m.findByTokenHashFn(ctx, hash)
}
func (m *mockSessionRepo) FindByTokenHashForUpdate(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error) {
	if m.findByTokenHashForUpdFn == nil {
		return nil, domain.ErrRefreshTokenNotFound
	}
	return m.findByTokenHashForUpdFn(ctx, hash)
}
func (m *mockSessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	if m.revokeFn == nil {
		return nil
	}
	return m.revokeFn(ctx, id)
}
func (m *mockSessionRepo) DeleteExpiredAndRevoked(ctx context.Context, userID uuid.UUID) error {
	if m.deleteExpiredFn == nil {
		return nil
	}
	return m.deleteExpiredFn(ctx, userID)
}
func (m *mockSessionRepo) RevokeAllByUserID(_ context.Context, _ uuid.UUID) error { return nil }

type mockHasher struct {
	verifyFn func(password, hash string) (bool, error)
}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Verify(password, hash string) (bool, error) {
	if m.verifyFn == nil {
		return password == "correct_password", nil
	}
	return m.verifyFn(password, hash)
}

type mockTokenProvider struct {
	generateAccessFn  func(ctx context.Context, claims value.AccessClaims) (string, int64, error)
	generateRefreshFn func(ctx context.Context) (string, value.TokenHash, error)
}

func (m *mockTokenProvider) GenerateAccessToken(ctx context.Context, claims value.AccessClaims) (string, int64, error) {
	if m.generateAccessFn == nil {
		return "access_token", 900, nil
	}
	return m.generateAccessFn(ctx, claims)
}
func (m *mockTokenProvider) ParseAccessToken(_ context.Context, _ string) (value.AccessClaims, error) {
	return value.AccessClaims{}, nil
}
func (m *mockTokenProvider) GenerateRefreshToken(ctx context.Context) (string, value.TokenHash, error) {
	if m.generateRefreshFn == nil {
		return "raw_refresh", value.TokenHash("hashed_refresh"), nil
	}
	return m.generateRefreshFn(ctx)
}
func (m *mockTokenProvider) GeneratePasswordResetToken(_ context.Context) (string, value.TokenHash, error) {
	return "raw_reset", value.TokenHash("hashed_reset"), nil
}

type mockTokenHasher struct{}

func (m *mockTokenHasher) Hash(raw string) value.TokenHash {
	return value.TokenHash("hashed_" + raw)
}

type mockClock struct{ now time.Time }

func (m *mockClock) Now() time.Time { return m.now }

type mockUUID struct{ id uuid.UUID }

func (m *mockUUID) New() uuid.UUID { return m.id }

type inlineTx struct{}

func (t *inlineTx) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ─── Stub repos для PermissionResolver ───────────────────────────────────────

type stubUserRoleRepo struct{}

func (r *stubUserRoleRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (r *stubUserRoleRepo) FindByUserIDForUpdate(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (r *stubUserRoleRepo) FindByUserIDs(_ context.Context, _ []uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (r *stubUserRoleRepo) Assign(_ context.Context, _ entity.UserRole) error      { return nil }
func (r *stubUserRoleRepo) Revoke(_ context.Context, _, _ uuid.UUID) error         { return nil }
func (r *stubUserRoleRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }
func (r *stubUserRoleRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}

type stubRoleRepo struct{}

func (r *stubRoleRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Role, error) {
	return nil, nil
}
func (r *stubRoleRepo) FindByCodes(_ context.Context, _ []string) ([]entity.Role, error) {
	return nil, nil
}
func (r *stubRoleRepo) FindByCode(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (r *stubRoleRepo) FindByCodeForUpdate(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (r *stubRoleRepo) FindAll(_ context.Context) ([]entity.Role, error)       { return nil, nil }
func (r *stubRoleRepo) ExistsByCode(_ context.Context, _ string) (bool, error) { return false, nil }
func (r *stubRoleRepo) Create(_ context.Context, _ entity.Role) error          { return nil }
func (r *stubRoleRepo) Delete(_ context.Context, _ uuid.UUID) error            { return nil }

type stubPermRepo struct{}

func (r *stubPermRepo) FindByID(_ context.Context, _ uuid.UUID) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (r *stubPermRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Permission, error) {
	return nil, nil
}
func (r *stubPermRepo) FindByCode(_ context.Context, _ string) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (r *stubPermRepo) FindByCodeForUpdate(_ context.Context, _ string) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (r *stubPermRepo) FindAll(_ context.Context) ([]entity.Permission, error) { return nil, nil }
func (r *stubPermRepo) ExistsByCode(_ context.Context, _ string) (bool, error) { return false, nil }
func (r *stubPermRepo) Create(_ context.Context, _ entity.Permission) error    { return nil }
func (r *stubPermRepo) Delete(_ context.Context, _ uuid.UUID) error            { return nil }

type stubRolePermRepo struct{}

func (r *stubRolePermRepo) FindByRoleID(_ context.Context, _ uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (r *stubRolePermRepo) FindByRoleIDs(_ context.Context, _ []uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (r *stubRolePermRepo) Assign(_ context.Context, _ entity.RolePermission) error { return nil }
func (r *stubRolePermRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error)  { return false, nil }
func (r *stubRolePermRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}
func (r *stubRolePermRepo) ExistsByPermissionID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}
func (r *stubRolePermRepo) Revoke(_ context.Context, _, _ uuid.UUID) error { return nil }

// ─── Builder ──────────────────────────────────────────────────────────────────

func makeResolver() *common.PermissionResolver {
	return common.NewPermissionResolver(
		&stubUserRoleRepo{},
		&stubRolePermRepo{},
		&stubRoleRepo{},
		&stubPermRepo{},
	)
}

type authBuilder struct {
	users    *mockUserRepo
	sessions *mockSessionRepo
	hasher   *mockHasher
	tokens   *mockTokenProvider
}

func newAuthBuilder() *authBuilder {
	return &authBuilder{
		users: &mockUserRepo{
			findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
				u := activeUser
				return &u, nil
			},
		},
		sessions: &mockSessionRepo{},
		hasher:   &mockHasher{},
		tokens:   &mockTokenProvider{},
	}
}

func (b *authBuilder) build() *authusecase.AuthService {
	return authusecase.NewAuthService(
		b.users,
		b.sessions,
		b.hasher,
		b.tokens,
		&mockTokenHasher{},
		&inlineTx{},
		&mockClock{now: testNow},
		&mockUUID{id: uuid.New()},
		makeResolver(),
		7*24*time.Hour,
	)
}

// ─── Login ────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	b := newAuthBuilder()
	b.hasher.verifyFn = func(_, _ string) (bool, error) { return true, nil }

	result, rawRefresh, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "admin@example.com",
		Password: "correct_password",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, rawRefresh)
	assert.Equal(t, testUserID, result.User.ID)
}

func TestLogin_InvalidEmail_MaskedAsInvalidCredentials(t *testing.T) {
	_, _, err := newAuthBuilder().build().Login(context.Background(), input.LoginInput{
		Email:    "not-an-email",
		Password: "pass",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_UserNotFound_MaskedAsInvalidCredentials(t *testing.T) {
	b := newAuthBuilder()
	b.users.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		return nil, domain.ErrUserNotFound
	}
	_, _, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "unknown@example.com",
		Password: "pass",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_InactiveUser(t *testing.T) {
	b := newAuthBuilder()
	b.users.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		u := activeUser
		u.IsActive = false
		return &u, nil
	}
	_, _, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "admin@example.com",
		Password: "correct_password",
	})
	assert.ErrorIs(t, err, domain.ErrUserInactive)
}

func TestLogin_WrongPassword(t *testing.T) {
	b := newAuthBuilder()
	b.hasher.verifyFn = func(_, _ string) (bool, error) { return false, nil }
	_, _, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "admin@example.com",
		Password: "wrong_password",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_InfraError_NotMaskedAsInvalidCredentials(t *testing.T) {
	b := newAuthBuilder()
	infraErr := errors.New("db connection lost")
	b.users.findByEmailFn = func(_ context.Context, _ value.Email) (*entity.User, error) {
		return nil, infraErr
	}
	_, _, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "admin@example.com",
		Password: "correct_password",
	})
	assert.False(t, errors.Is(err, domain.ErrInvalidCredentials))
	assert.ErrorContains(t, err, "db connection lost")
}

func TestLogin_SessionSaved(t *testing.T) {
	var sessionSaved bool
	b := newAuthBuilder()
	b.hasher.verifyFn = func(_, _ string) (bool, error) { return true, nil }
	b.sessions.saveFn = func(_ context.Context, _ entity.RefreshSession) error {
		sessionSaved = true
		return nil
	}
	_, _, err := b.build().Login(context.Background(), input.LoginInput{
		Email:    "admin@example.com",
		Password: "correct_password",
	})
	require.NoError(t, err)
	assert.True(t, sessionSaved)
}

// ─── Refresh ──────────────────────────────────────────────────────────────────

func makeActiveSession(userID uuid.UUID, now time.Time) entity.RefreshSession {
	s, _ := entity.NewRefreshSession(uuid.New(), userID, value.TokenHash("hashed_raw_token"), now, time.Hour)
	return s
}

func TestRefresh_Success(t *testing.T) {
	session := makeActiveSession(testUserID, testNow)
	b := newAuthBuilder()
	b.users.findByIDFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		u := activeUser
		return &u, nil
	}
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}

	result, rawRefresh, err := b.build().Refresh(context.Background(), input.RefreshInput{
		RawRefreshToken: "raw_token",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, rawRefresh)
}

func TestRefresh_SessionNotFound(t *testing.T) {
	b := newAuthBuilder()
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return nil, domain.ErrRefreshTokenNotFound
	}
	_, _, err := b.build().Refresh(context.Background(), input.RefreshInput{RawRefreshToken: "bad"})
	assert.Error(t, err)
}

func TestRefresh_SessionRevoked(t *testing.T) {
	session := makeActiveSession(testUserID, testNow)
	revokedAt := testNow
	session.RevokedAt = &revokedAt

	b := newAuthBuilder()
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	_, _, err := b.build().Refresh(context.Background(), input.RefreshInput{RawRefreshToken: "raw_token"})
	assert.ErrorIs(t, err, domain.ErrRefreshTokenRevoked)
}

func TestRefresh_SessionExpired(t *testing.T) {
	session := makeActiveSession(testUserID, testNow) // ExpiresAt = testNow + 1h
	expiredNow := testNow.Add(2 * time.Hour)

	b := newAuthBuilder()
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	b.users.findByIDFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		u := activeUser
		return &u, nil
	}

	svc := authusecase.NewAuthService(
		b.users, b.sessions, b.hasher, b.tokens, &mockTokenHasher{},
		&inlineTx{}, &mockClock{now: expiredNow}, &mockUUID{id: uuid.New()},
		makeResolver(), 7*24*time.Hour,
	)
	_, _, err := svc.Refresh(context.Background(), input.RefreshInput{RawRefreshToken: "raw_token"})
	assert.ErrorIs(t, err, domain.ErrRefreshTokenExpired)
}

func TestRefresh_InactiveUser(t *testing.T) {
	session := makeActiveSession(testUserID, testNow)
	b := newAuthBuilder()
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	b.users.findByIDFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		u := activeUser
		u.IsActive = false
		return &u, nil
	}
	_, _, err := b.build().Refresh(context.Background(), input.RefreshInput{RawRefreshToken: "raw_token"})
	assert.ErrorIs(t, err, domain.ErrUserInactive)
}

func TestRefresh_OldSessionRevoked_NewSessionSaved(t *testing.T) {
	session := makeActiveSession(testUserID, testNow)
	var oldRevoked, newSaved bool

	b := newAuthBuilder()
	b.sessions.findByTokenHashForUpdFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	b.users.findByIDFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		u := activeUser
		return &u, nil
	}
	b.sessions.revokeFn = func(_ context.Context, _ uuid.UUID) error {
		oldRevoked = true
		return nil
	}
	b.sessions.saveFn = func(_ context.Context, _ entity.RefreshSession) error {
		newSaved = true
		return nil
	}

	_, _, err := b.build().Refresh(context.Background(), input.RefreshInput{RawRefreshToken: "raw_token"})
	require.NoError(t, err)
	assert.True(t, oldRevoked, "old session must be revoked")
	assert.True(t, newSaved, "new session must be saved")
}

// ─── Logout ───────────────────────────────────────────────────────────────────

func TestLogout_ActiveSession_Revoked(t *testing.T) {
	session := makeActiveSession(testUserID, testNow)
	var revoked bool

	b := newAuthBuilder()
	b.sessions.findByTokenHashFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	b.sessions.revokeFn = func(_ context.Context, _ uuid.UUID) error {
		revoked = true
		return nil
	}

	err := b.build().Logout(context.Background(), input.LogoutInput{RawRefreshToken: "raw_token"})
	require.NoError(t, err)
	assert.True(t, revoked)
}

func TestLogout_SessionNotFound_SilentSuccess(t *testing.T) {
	b := newAuthBuilder()
	b.sessions.findByTokenHashFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return nil, domain.ErrRefreshTokenNotFound
	}
	err := b.build().Logout(context.Background(), input.LogoutInput{RawRefreshToken: "bad_token"})
	assert.NoError(t, err)
}

func TestLogout_RevokedSession_SilentSuccess(t *testing.T) {
	revokedAt := testNow
	session := makeActiveSession(testUserID, testNow)
	session.RevokedAt = &revokedAt

	b := newAuthBuilder()
	b.sessions.findByTokenHashFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return &session, nil
	}
	b.sessions.revokeFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("must not call Revoke for already revoked session")
		return nil
	}

	err := b.build().Logout(context.Background(), input.LogoutInput{RawRefreshToken: "raw_token"})
	assert.NoError(t, err)
}

func TestLogout_InfraError_Propagated(t *testing.T) {
	infraErr := errors.New("db error")
	b := newAuthBuilder()
	b.sessions.findByTokenHashFn = func(_ context.Context, _ value.TokenHash) (*entity.RefreshSession, error) {
		return nil, infraErr
	}
	err := b.build().Logout(context.Background(), input.LogoutInput{RawRefreshToken: "raw_token"})
	assert.ErrorContains(t, err, "db error")
}
