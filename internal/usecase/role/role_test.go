package role_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	roleusecase "github.com/AlbinaKonovalova/auth-service/internal/usecase/role"
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockRoleRepo struct {
	existsByCodeFn     func(ctx context.Context, code string) (bool, error)
	createFn           func(ctx context.Context, r entity.Role) error
	findByCodeForUpdFn func(ctx context.Context, code string) (*entity.Role, error)
	deleteFn           func(ctx context.Context, id uuid.UUID) error
}

func (m *mockRoleRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	if m.existsByCodeFn == nil {
		return false, nil
	}
	return m.existsByCodeFn(ctx, code)
}
func (m *mockRoleRepo) Create(ctx context.Context, r entity.Role) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(ctx, r)
}
func (m *mockRoleRepo) FindByCodeForUpdate(ctx context.Context, code string) (*entity.Role, error) {
	if m.findByCodeForUpdFn == nil {
		return &entity.Role{ID: uuid.New(), Code: code, Name: "Test"}, nil
	}
	return m.findByCodeForUpdFn(ctx, code)
}
func (m *mockRoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(ctx, id)
}
func (m *mockRoleRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCodes(_ context.Context, _ []string) ([]entity.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCode(_ context.Context, _ string) (*entity.Role, error) {
	return nil, domain.ErrRoleNotFound
}
func (m *mockRoleRepo) FindAll(_ context.Context) ([]entity.Role, error) { return []entity.Role{}, nil }

type mockUserRoleRepo struct {
	existsByRoleIDFn func(ctx context.Context, roleID uuid.UUID) (bool, error)
}

func (m *mockUserRoleRepo) ExistsByRoleID(ctx context.Context, roleID uuid.UUID) (bool, error) {
	if m.existsByRoleIDFn == nil {
		return false, nil
	}
	return m.existsByRoleIDFn(ctx, roleID)
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
func (m *mockUserRoleRepo) Assign(_ context.Context, _ entity.UserRole) error      { return nil }
func (m *mockUserRoleRepo) Revoke(_ context.Context, _, _ uuid.UUID) error         { return nil }
func (m *mockUserRoleRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }

type mockRolePermRepo struct {
	existsByRoleIDFn func(ctx context.Context, roleID uuid.UUID) (bool, error)
}

func (m *mockRolePermRepo) ExistsByRoleID(ctx context.Context, roleID uuid.UUID) (bool, error) {
	if m.existsByRoleIDFn == nil {
		return false, nil
	}
	return m.existsByRoleIDFn(ctx, roleID)
}
func (m *mockRolePermRepo) FindByRoleID(_ context.Context, _ uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) FindByRoleIDs(_ context.Context, _ []uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) Assign(_ context.Context, _ entity.RolePermission) error { return nil }
func (m *mockRolePermRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error)  { return false, nil }
func (m *mockRolePermRepo) ExistsByPermissionID(_ context.Context, _ uuid.UUID) (bool, error) {
	return false, nil
}
func (m *mockRolePermRepo) Revoke(_ context.Context, _, _ uuid.UUID) error { return nil }

type mockUUID struct{}

func (m *mockUUID) New() uuid.UUID { return uuid.New() }

type inlineTx struct{}

func (t *inlineTx) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ─── Builder ──────────────────────────────────────────────────────────────────

type roleBuilder struct {
	roleRepo     *mockRoleRepo
	userRoleRepo *mockUserRoleRepo
	rolePermRepo *mockRolePermRepo
}

func newRoleBuilder() *roleBuilder {
	return &roleBuilder{
		roleRepo:     &mockRoleRepo{},
		userRoleRepo: &mockUserRoleRepo{},
		rolePermRepo: &mockRolePermRepo{},
	}
}

func (b *roleBuilder) build() *roleusecase.RoleService {
	return roleusecase.NewRoleService(b.roleRepo, b.userRoleRepo, b.rolePermRepo, &mockUUID{}, &inlineTx{})
}

// ─── CreateRole ───────────────────────────────────────────────────────────────

func TestCreateRole_Success(t *testing.T) {
	b := newRoleBuilder()
	var created bool
	b.roleRepo.createFn = func(_ context.Context, _ entity.Role) error {
		created = true
		return nil
	}

	svc := b.build()
	view, err := svc.CreateRole(context.Background(), input.CreateRoleInput{
		Code: "manager", Name: "Manager", Description: "Manages things",
	})

	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, "manager", view.Code)
	assert.Equal(t, "Manager", view.Name)
}

func TestCreateRole_InvalidCode(t *testing.T) {
	svc := newRoleBuilder().build()
	_, err := svc.CreateRole(context.Background(), input.CreateRoleInput{
		Code: "", Name: "Manager",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestCreateRole_EmptyName(t *testing.T) {
	svc := newRoleBuilder().build()
	_, err := svc.CreateRole(context.Background(), input.CreateRoleInput{
		Code: "manager", Name: "",
	})
	assert.ErrorIs(t, err, domain.ErrRoleNameEmpty)
}

func TestCreateRole_DuplicateCode(t *testing.T) {
	b := newRoleBuilder()
	b.roleRepo.existsByCodeFn = func(_ context.Context, _ string) (bool, error) {
		return true, nil
	}
	svc := b.build()
	_, err := svc.CreateRole(context.Background(), input.CreateRoleInput{
		Code: "admin", Name: "Admin",
	})
	assert.ErrorIs(t, err, domain.ErrDuplicateRoleCode)
}

func TestCreateRole_CodeNormalized(t *testing.T) {
	b := newRoleBuilder()
	var createdCode string
	b.roleRepo.createFn = func(_ context.Context, r entity.Role) error {
		createdCode = r.Code
		return nil
	}

	svc := b.build()
	_, err := svc.CreateRole(context.Background(), input.CreateRoleInput{
		Code: "MANAGER", Name: "Manager",
	})
	require.NoError(t, err)
	assert.Equal(t, "manager", createdCode) // нормализован к lowercase
}

// ─── DeleteRole ───────────────────────────────────────────────────────────────

func TestDeleteRole_Success(t *testing.T) {
	b := newRoleBuilder()
	var deleted bool
	b.roleRepo.deleteFn = func(_ context.Context, _ uuid.UUID) error {
		deleted = true
		return nil
	}

	svc := b.build()
	err := svc.DeleteRole(context.Background(), "admin")
	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteRole_InvalidCode(t *testing.T) {
	svc := newRoleBuilder().build()
	err := svc.DeleteRole(context.Background(), "")
	assert.ErrorIs(t, err, domain.ErrInvalidRoleCode)
}

func TestDeleteRole_NotFound(t *testing.T) {
	b := newRoleBuilder()
	b.roleRepo.findByCodeForUpdFn = func(_ context.Context, _ string) (*entity.Role, error) {
		return nil, domain.ErrRoleNotFound
	}
	svc := b.build()
	err := svc.DeleteRole(context.Background(), "ghost")
	assert.Error(t, err)
}

func TestDeleteRole_HasUsers_Blocked(t *testing.T) {
	b := newRoleBuilder()
	b.userRoleRepo.existsByRoleIDFn = func(_ context.Context, _ uuid.UUID) (bool, error) {
		return true, nil
	}
	svc := b.build()
	err := svc.DeleteRole(context.Background(), "admin")
	assert.ErrorIs(t, err, domain.ErrRoleInUse)
}

func TestDeleteRole_HasPermissions_Blocked(t *testing.T) {
	b := newRoleBuilder()
	b.rolePermRepo.existsByRoleIDFn = func(_ context.Context, _ uuid.UUID) (bool, error) {
		return true, nil
	}
	svc := b.build()
	err := svc.DeleteRole(context.Background(), "admin")
	assert.ErrorIs(t, err, domain.ErrRoleInUse)
}
