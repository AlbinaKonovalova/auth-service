package access_test

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
	accessusecase "github.com/AlbinaKonovalova/auth-service/internal/usecase/access"
)

var (
	testNow    = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	testUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	roleID1    = uuid.MustParse("00000000-0000-0000-0000-000000000010")
	roleID2    = uuid.MustParse("00000000-0000-0000-0000-000000000011")
	permID1    = uuid.MustParse("00000000-0000-0000-0000-000000000020")

	roleAdmin = entity.Role{ID: roleID1, Code: "admin", Name: "Admin"}
	permRead  = entity.Permission{ID: permID1, Code: "users.read", Description: "Read"}
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	findByIDFn func(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDFn == nil {
		return &entity.User{ID: id, IsActive: true}, nil
	}
	return m.findByIDFn(ctx, id)
}
func (m *mockUserRepo) FindByIDForUpdate(_ context.Context, id uuid.UUID) (*entity.User, error) {
	return &entity.User{ID: id, IsActive: true}, nil
}
func (m *mockUserRepo) FindByEmail(_ context.Context, _ value.Email) (*entity.User, error) {
	return nil, domain.ErrUserNotFound
}
func (m *mockUserRepo) Create(_ context.Context, _ entity.User) error                     { return nil }
func (m *mockUserRepo) UpdatePasswordHash(_ context.Context, _ uuid.UUID, _ string) error { return nil }
func (m *mockUserRepo) List(_ context.Context, _ dto.UserListFilters) ([]entity.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Activate(_ context.Context, _ uuid.UUID) error   { return nil }
func (m *mockUserRepo) Deactivate(_ context.Context, _ uuid.UUID) error { return nil }

type mockUserRoleRepo struct {
	findByUserIDForUpdFn func(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error)
	existsFn             func(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	assignFn             func(ctx context.Context, ur entity.UserRole) error
	revokeFn             func(ctx context.Context, userID, roleID uuid.UUID) error
}

func (m *mockUserRoleRepo) FindByUserIDForUpdate(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error) {
	if m.findByUserIDForUpdFn == nil {
		return nil, nil
	}
	return m.findByUserIDForUpdFn(ctx, userID)
}
func (m *mockUserRoleRepo) Exists(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	if m.existsFn == nil {
		return false, nil
	}
	return m.existsFn(ctx, userID, roleID)
}
func (m *mockUserRoleRepo) Assign(ctx context.Context, ur entity.UserRole) error {
	if m.assignFn == nil {
		return nil
	}
	return m.assignFn(ctx, ur)
}
func (m *mockUserRoleRepo) Revoke(ctx context.Context, userID, roleID uuid.UUID) error {
	if m.revokeFn == nil {
		return nil
	}
	return m.revokeFn(ctx, userID, roleID)
}
func (m *mockUserRoleRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (m *mockUserRoleRepo) FindByUserIDs(_ context.Context, _ []uuid.UUID) ([]entity.UserRole, error) {
	return nil, nil
}
func (m *mockUserRoleRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}

type mockRoleRepo struct {
	findByCodeFn func(ctx context.Context, code string) (*entity.Role, error)
}

func (m *mockRoleRepo) FindByCode(ctx context.Context, code string) (*entity.Role, error) {
	if m.findByCodeFn == nil {
		return &roleAdmin, nil
	}
	return m.findByCodeFn(ctx, code)
}
func (m *mockRoleRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCodes(_ context.Context, _ []string) ([]entity.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCodeForUpdate(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (m *mockRoleRepo) FindAll(_ context.Context) ([]entity.Role, error)       { return nil, nil }
func (m *mockRoleRepo) ExistsByCode(_ context.Context, _ string) (bool, error) { return false, nil }
func (m *mockRoleRepo) Create(_ context.Context, _ entity.Role) error          { return nil }
func (m *mockRoleRepo) Delete(_ context.Context, _ uuid.UUID) error            { return nil }

type mockRolePermRepo struct {
	assignFn func(ctx context.Context, rp entity.RolePermission) error
	existsFn func(ctx context.Context, roleID, permID uuid.UUID) (bool, error)
	revokeFn func(ctx context.Context, roleID, permID uuid.UUID) error
}

func (m *mockRolePermRepo) Assign(ctx context.Context, rp entity.RolePermission) error {
	if m.assignFn == nil {
		return nil
	}
	return m.assignFn(ctx, rp)
}
func (m *mockRolePermRepo) Exists(ctx context.Context, roleID, permID uuid.UUID) (bool, error) {
	if m.existsFn == nil {
		return true, nil
	}
	return m.existsFn(ctx, roleID, permID)
}
func (m *mockRolePermRepo) Revoke(ctx context.Context, roleID, permID uuid.UUID) error {
	if m.revokeFn == nil {
		return nil
	}
	return m.revokeFn(ctx, roleID, permID)
}
func (m *mockRolePermRepo) FindByRoleID(_ context.Context, _ uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) FindByRoleIDs(_ context.Context, _ []uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}
func (m *mockRolePermRepo) ExistsByPermissionID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}

type mockPermRepo struct {
	findByCodeFn func(ctx context.Context, code string) (*entity.Permission, error)
}

func (m *mockPermRepo) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	if m.findByCodeFn == nil {
		return &permRead, nil
	}
	return m.findByCodeFn(ctx, code)
}
func (m *mockPermRepo) FindByID(_ context.Context, _ uuid.UUID) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (m *mockPermRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Permission, error) {
	return nil, nil
}
func (m *mockPermRepo) FindByCodeForUpdate(_ context.Context, _ string) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (m *mockPermRepo) FindAll(_ context.Context) ([]entity.Permission, error) { return nil, nil }
func (m *mockPermRepo) ExistsByCode(_ context.Context, _ string) (bool, error) { return false, nil }
func (m *mockPermRepo) Create(_ context.Context, _ entity.Permission) error    { return nil }
func (m *mockPermRepo) Delete(_ context.Context, _ uuid.UUID) error            { return nil }

type mockClock struct{}

func (m *mockClock) Now() time.Time { return testNow }

type inlineTx struct{}

func (t *inlineTx) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ─── Builder ──────────────────────────────────────────────────────────────────

type accessBuilder struct {
	userRepo     *mockUserRepo
	userRoleRepo *mockUserRoleRepo
	roleRepo     *mockRoleRepo
	rolePermRepo *mockRolePermRepo
	permRepo     *mockPermRepo
}

func newAccessBuilder() *accessBuilder {
	return &accessBuilder{
		userRepo:     &mockUserRepo{},
		userRoleRepo: &mockUserRoleRepo{},
		roleRepo:     &mockRoleRepo{},
		rolePermRepo: &mockRolePermRepo{},
		permRepo:     &mockPermRepo{},
	}
}

func (b *accessBuilder) build() *accessusecase.AccessService {
	return accessusecase.NewAccessService(
		b.userRepo, b.userRoleRepo, b.roleRepo, b.rolePermRepo, b.permRepo,
		&mockClock{}, &inlineTx{},
	)
}

// ─── AssignRole ───────────────────────────────────────────────────────────────

func TestAssignRole_Success(t *testing.T) {
	var assigned bool
	b := newAccessBuilder()
	b.userRoleRepo.assignFn = func(_ context.Context, _ entity.UserRole) error {
		assigned = true
		return nil
	}
	err := b.build().AssignRole(context.Background(), input.AssignRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	require.NoError(t, err)
	assert.True(t, assigned)
}

func TestAssignRole_InvalidRoleCode(t *testing.T) {
	err := newAccessBuilder().build().AssignRole(context.Background(), input.AssignRoleInput{
		UserID: testUserID, RoleCode: "",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestAssignRole_UserNotFound(t *testing.T) {
	b := newAccessBuilder()
	b.userRepo.findByIDFn = func(_ context.Context, _ uuid.UUID) (*entity.User, error) {
		return nil, domain.ErrUserNotFound
	}
	err := b.build().AssignRole(context.Background(), input.AssignRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	assert.Error(t, err)
}

func TestAssignRole_RoleNotFound(t *testing.T) {
	b := newAccessBuilder()
	b.roleRepo.findByCodeFn = func(_ context.Context, _ string) (*entity.Role, error) {
		return nil, domain.ErrRoleNotFound
	}
	err := b.build().AssignRole(context.Background(), input.AssignRoleInput{
		UserID: testUserID, RoleCode: "ghost",
	})
	assert.Error(t, err)
}

func TestAssignRole_AlreadyAssigned_Idempotent(t *testing.T) {
	b := newAccessBuilder()
	b.userRoleRepo.existsFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return true, nil }
	b.userRoleRepo.assignFn = func(_ context.Context, _ entity.UserRole) error {
		t.Error("Assign must not be called if role already assigned")
		return nil
	}
	err := b.build().AssignRole(context.Background(), input.AssignRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	assert.NoError(t, err)
}

// ─── RevokeRole ───────────────────────────────────────────────────────────────

func TestRevokeRole_Success(t *testing.T) {
	var revoked bool
	b := newAccessBuilder()
	b.userRoleRepo.findByUserIDForUpdFn = func(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
		return []entity.UserRole{
			{UserID: testUserID, RoleID: roleID1},
			{UserID: testUserID, RoleID: roleID2},
		}, nil
	}
	b.userRoleRepo.revokeFn = func(_ context.Context, _, _ uuid.UUID) error {
		revoked = true
		return nil
	}
	err := b.build().RevokeRole(context.Background(), input.RevokeRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	require.NoError(t, err)
	assert.True(t, revoked)
}

func TestRevokeRole_InvalidRoleCode(t *testing.T) {
	err := newAccessBuilder().build().RevokeRole(context.Background(), input.RevokeRoleInput{
		UserID: testUserID, RoleCode: "",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestRevokeRole_RoleNotAssigned(t *testing.T) {
	b := newAccessBuilder()
	b.userRoleRepo.findByUserIDForUpdFn = func(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
		return []entity.UserRole{{UserID: testUserID, RoleID: roleID2}}, nil // только manager
	}
	// ищем admin (roleID1) — его нет у пользователя
	err := b.build().RevokeRole(context.Background(), input.RevokeRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	assert.ErrorIs(t, err, domain.ErrUserRoleNotFound)
}

func TestRevokeRole_LastRole_Blocked(t *testing.T) {
	b := newAccessBuilder()
	b.userRoleRepo.findByUserIDForUpdFn = func(_ context.Context, _ uuid.UUID) ([]entity.UserRole, error) {
		return []entity.UserRole{{UserID: testUserID, RoleID: roleID1}}, nil // единственная роль
	}
	err := b.build().RevokeRole(context.Background(), input.RevokeRoleInput{
		UserID: testUserID, RoleCode: "admin",
	})
	assert.ErrorIs(t, err, domain.ErrCannotRevokeLastRole)
}

// ─── AssignPermission ─────────────────────────────────────────────────────────

func TestAssignPermission_Success(t *testing.T) {
	var assigned bool
	b := newAccessBuilder()
	b.rolePermRepo.assignFn = func(_ context.Context, _ entity.RolePermission) error {
		assigned = true
		return nil
	}
	err := b.build().AssignPermission(context.Background(), input.AssignPermissionInput{
		RoleCode: "admin", PermissionCode: "users.read",
	})
	require.NoError(t, err)
	assert.True(t, assigned)
}

func TestAssignPermission_InvalidRoleCode(t *testing.T) {
	err := newAccessBuilder().build().AssignPermission(context.Background(), input.AssignPermissionInput{
		RoleCode: "", PermissionCode: "users.read",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestAssignPermission_InvalidPermissionCode(t *testing.T) {
	err := newAccessBuilder().build().AssignPermission(context.Background(), input.AssignPermissionInput{
		RoleCode: "admin", PermissionCode: "",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidPermissionCode)
}

func TestAssignPermission_RoleNotFound(t *testing.T) {
	b := newAccessBuilder()
	b.roleRepo.findByCodeFn = func(_ context.Context, _ string) (*entity.Role, error) {
		return nil, domain.ErrRoleNotFound
	}
	err := b.build().AssignPermission(context.Background(), input.AssignPermissionInput{
		RoleCode: "ghost", PermissionCode: "users.read",
	})
	assert.Error(t, err)
}

func TestAssignPermission_PermissionNotFound(t *testing.T) {
	b := newAccessBuilder()
	b.permRepo.findByCodeFn = func(_ context.Context, _ string) (*entity.Permission, error) {
		return nil, domain.ErrPermissionNotFound
	}
	err := b.build().AssignPermission(context.Background(), input.AssignPermissionInput{
		RoleCode: "admin", PermissionCode: "ghost.perm",
	})
	assert.Error(t, err)
}

// ─── RevokePermission ─────────────────────────────────────────────────────────

func TestRevokePermission_Success(t *testing.T) {
	var revoked bool
	b := newAccessBuilder()
	b.rolePermRepo.existsFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return true, nil }
	b.rolePermRepo.revokeFn = func(_ context.Context, _, _ uuid.UUID) error {
		revoked = true
		return nil
	}
	err := b.build().RevokePermission(context.Background(), input.RevokePermissionInput{
		RoleCode: "admin", PermissionCode: "users.read",
	})
	require.NoError(t, err)
	assert.True(t, revoked)
}

func TestRevokePermission_InvalidRoleCode(t *testing.T) {
	err := newAccessBuilder().build().RevokePermission(context.Background(), input.RevokePermissionInput{
		RoleCode: "", PermissionCode: "users.read",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestRevokePermission_InvalidPermissionCode(t *testing.T) {
	err := newAccessBuilder().build().RevokePermission(context.Background(), input.RevokePermissionInput{
		RoleCode: "admin", PermissionCode: "",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidPermissionCode)
}

func TestRevokePermission_NotAssigned(t *testing.T) {
	b := newAccessBuilder()
	b.rolePermRepo.existsFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }
	err := b.build().RevokePermission(context.Background(), input.RevokePermissionInput{
		RoleCode: "admin", PermissionCode: "users.read",
	})
	assert.ErrorIs(t, err, domain.ErrRolePermissionNotFound)
}
