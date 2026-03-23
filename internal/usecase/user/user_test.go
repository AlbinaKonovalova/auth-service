package user_test

import (
	"context"
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
	userusecase "github.com/AlbinaKonovalova/auth-service/internal/usecase/user"
)

// ─── Fixtures ─────────────────────────────────────────────────────────────────

var (
	testNow     = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	testUserID  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	adminRoleID = uuid.MustParse("00000000-0000-0000-0000-000000000010")
	adminRole   = entity.Role{ID: adminRoleID, Code: "admin", Name: "Admin"}
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	createFn            func(ctx context.Context, u entity.User) error
	findByIDForUpdateFn func(ctx context.Context, id uuid.UUID) (*entity.User, error)
	activateFn          func(ctx context.Context, id uuid.UUID) error
	deactivateFn        func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepo) Create(ctx context.Context, u entity.User) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(ctx, u)
}
func (m *mockUserRepo) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDForUpdateFn == nil {
		return &entity.User{ID: id, IsActive: true}, nil
	}
	return m.findByIDForUpdateFn(ctx, id)
}
func (m *mockUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	return &entity.User{ID: id, IsActive: true}, nil
}
func (m *mockUserRepo) FindByEmail(_ context.Context, _ value.Email) (*entity.User, error) {
	return nil, domain.ErrUserNotFound
}
func (m *mockUserRepo) UpdatePasswordHash(_ context.Context, _ uuid.UUID, _ string) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _ dto.UserListFilters) ([]entity.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Activate(ctx context.Context, id uuid.UUID) error {
	if m.activateFn == nil {
		return nil
	}
	return m.activateFn(ctx, id)
}
func (m *mockUserRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	if m.deactivateFn == nil {
		return nil
	}
	return m.deactivateFn(ctx, id)
}

type mockRoleRepo struct {
	findByCodesFn func(ctx context.Context, codes []string) ([]entity.Role, error)
}

func (m *mockRoleRepo) FindByCodes(ctx context.Context, codes []string) ([]entity.Role, error) {
	if m.findByCodesFn == nil {
		return []entity.Role{adminRole}, nil
	}
	return m.findByCodesFn(ctx, codes)
}
func (m *mockRoleRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCode(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (m *mockRoleRepo) FindByCodeForUpdate(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (m *mockRoleRepo) FindAll(_ context.Context) ([]entity.Role, error)       { return nil, nil }
func (m *mockRoleRepo) ExistsByCode(_ context.Context, _ string) (bool, error) { return false, nil }
func (m *mockRoleRepo) Create(_ context.Context, _ entity.Role) error          { return nil }
func (m *mockRoleRepo) Delete(_ context.Context, _ uuid.UUID) error            { return nil }

type mockUserRoleRepo struct {
	assignFn func(ctx context.Context, ur entity.UserRole) error
}

func (m *mockUserRoleRepo) Assign(ctx context.Context, ur entity.UserRole) error {
	if m.assignFn == nil {
		return nil
	}
	return m.assignFn(ctx, ur)
}
func (m *mockUserRoleRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (m *mockUserRoleRepo) FindByUserIDForUpdate(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (m *mockUserRoleRepo) FindByUserIDs(_ context.Context, _ []uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (m *mockUserRoleRepo) Revoke(_ context.Context, _, _ uuid.UUID) error         { return nil }
func (m *mockUserRoleRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }
func (m *mockUserRoleRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
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

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Verify(_, _ string) (bool, error)     { return true, nil }

type mockClock struct{}

func (m *mockClock) Now() time.Time { return testNow }

type mockUUID struct{}

func (m *mockUUID) New() uuid.UUID { return testUserID }

type inlineTx struct{}

func (t *inlineTx) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ─── Builder ──────────────────────────────────────────────────────────────────

type userBuilder struct {
	users     *mockUserRepo
	roles     *mockRoleRepo
	userRoles *mockUserRoleRepo
	sessions  *mockSessionRepo
}

func newUserBuilder() *userBuilder {
	return &userBuilder{
		users:     &mockUserRepo{},
		roles:     &mockRoleRepo{},
		userRoles: &mockUserRoleRepo{},
		sessions:  &mockSessionRepo{},
	}
}

func (b *userBuilder) build() *userusecase.UserService {
	return userusecase.NewUserService(
		b.users, b.roles, b.userRoles, b.sessions,
		&mockHasher{}, &inlineTx{}, &mockClock{}, &mockUUID{},
	)
}

// ─── CreateUser ───────────────────────────────────────────────────────────────

func TestCreateUser_Success(t *testing.T) {
	b := newUserBuilder()
	var userCreated bool
	b.users.createFn = func(_ context.Context, _ entity.User) error {
		userCreated = true
		return nil
	}

	svc := b.build()
	view, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "new@example.com",
		Password: "ValidPass1!",
		Roles:    []string{"admin"},
	})

	require.NoError(t, err)
	assert.True(t, userCreated)
	assert.True(t, view.IsActive)
	assert.Equal(t, []string{"admin"}, view.Roles)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	svc := newUserBuilder().build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "not-an-email",
		Password: "ValidPass1!",
		Roles:    []string{"admin"},
	})
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)
}

func TestCreateUser_InvalidPassword(t *testing.T) {
	svc := newUserBuilder().build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "user@example.com",
		Password: "short",
		Roles:    []string{"admin"},
	})
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}

func TestCreateUser_NoRoles(t *testing.T) {
	svc := newUserBuilder().build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "user@example.com",
		Password: "ValidPass1!",
		Roles:    []string{},
	})
	assert.ErrorIs(t, err, domain.ErrUserMustHaveRole)
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	b := newUserBuilder()
	b.roles.findByCodesFn = func(_ context.Context, _ []string) ([]entity.Role, error) {
		return []entity.Role{}, nil // роль не найдена
	}
	svc := b.build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "user@example.com",
		Password: "ValidPass1!",
		Roles:    []string{"nonexistent"},
	})
	assert.ErrorIs(t, err, domain.ErrRoleNotFound)
}

func TestCreateUser_DuplicateRoleCode(t *testing.T) {
	svc := newUserBuilder().build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "user@example.com",
		Password: "ValidPass1!",
		Roles:    []string{"admin", "ADMIN"}, // дубликат после нормализации
	})
	assert.ErrorIs(t, err, domain.ErrDuplicateRoleCode)
}

func TestCreateUser_EmailAlreadyTaken(t *testing.T) {
	b := newUserBuilder()
	b.users.createFn = func(_ context.Context, _ entity.User) error {
		return domain.ErrEmailAlreadyTaken
	}
	svc := b.build()
	_, err := svc.CreateUser(context.Background(), input.CreateUserInput{
		Email:    "existing@example.com",
		Password: "ValidPass1!",
		Roles:    []string{"admin"},
	})
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyTaken)
}

// ─── ActivateUser ─────────────────────────────────────────────────────────────

func TestActivateUser_Success(t *testing.T) {
	b := newUserBuilder()
	var activated bool
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return &entity.User{ID: testUserID, IsActive: false}, nil
	}
	b.users.activateFn = func(_ context.Context, _ uuid.UUID) error {
		activated = true
		return nil
	}

	svc := b.build()
	err := svc.ActivateUser(context.Background(), testUserID)
	require.NoError(t, err)
	assert.True(t, activated)
}

func TestActivateUser_AlreadyActive_Idempotent(t *testing.T) {
	b := newUserBuilder()
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return &entity.User{ID: testUserID, IsActive: true}, nil // уже активен
	}
	b.users.activateFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("Activate must not be called if user is already active")
		return nil
	}

	svc := b.build()
	err := svc.ActivateUser(context.Background(), testUserID)
	assert.NoError(t, err)
}

func TestActivateUser_UserNotFound(t *testing.T) {
	b := newUserBuilder()
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return nil, domain.ErrUserNotFound
	}
	svc := b.build()
	err := svc.ActivateUser(context.Background(), testUserID)
	assert.Error(t, err)
}

// ─── DeactivateUser ───────────────────────────────────────────────────────────

func TestDeactivateUser_Success_RevokesSession(t *testing.T) {
	var deactivated, sessionsRevoked bool

	b := newUserBuilder()
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return &entity.User{ID: testUserID, IsActive: true}, nil
	}
	b.users.deactivateFn = func(_ context.Context, _ uuid.UUID) error {
		deactivated = true
		return nil
	}
	b.sessions.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		sessionsRevoked = true
		return nil
	}

	svc := b.build()
	err := svc.DeactivateUser(context.Background(), testUserID)
	require.NoError(t, err)
	assert.True(t, deactivated)
	assert.True(t, sessionsRevoked)
}

func TestDeactivateUser_AlreadyInactive_SessionsStillRevoked(t *testing.T) {
	var sessionsRevoked bool

	b := newUserBuilder()
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return &entity.User{ID: testUserID, IsActive: false}, nil // уже неактивен
	}
	b.users.deactivateFn = func(_ context.Context, _ uuid.UUID) error {
		t.Error("Deactivate must not be called if already inactive")
		return nil
	}
	b.sessions.revokeAllFn = func(_ context.Context, _ uuid.UUID) error {
		sessionsRevoked = true
		return nil
	}

	svc := b.build()
	err := svc.DeactivateUser(context.Background(), testUserID)
	require.NoError(t, err)
	// Сессии должны отзываться даже если пользователь уже неактивен
	assert.True(t, sessionsRevoked)
}

func TestDeactivateUser_UserNotFound(t *testing.T) {
	b := newUserBuilder()
	b.users.findByIDForUpdateFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return nil, domain.ErrUserNotFound
	}
	svc := b.build()
	err := svc.DeactivateUser(context.Background(), testUserID)
	assert.Error(t, err)
}
