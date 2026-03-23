package permission_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	permusecase "github.com/AlbinaKonovalova/auth-service/internal/usecase/permission"
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockPermRepo struct {
	existsByCodeFn     func(ctx context.Context, code string) (bool, error)
	createFn           func(ctx context.Context, p entity.Permission) error
	findByCodeForUpdFn func(ctx context.Context, code string) (*entity.Permission, error)
	deleteFn           func(ctx context.Context, id uuid.UUID) error
}

func (m *mockPermRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	if m.existsByCodeFn == nil {
		return false, nil
	}
	return m.existsByCodeFn(ctx, code)
}
func (m *mockPermRepo) Create(ctx context.Context, p entity.Permission) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(ctx, p)
}
func (m *mockPermRepo) FindByCodeForUpdate(ctx context.Context, code string) (*entity.Permission, error) {
	if m.findByCodeForUpdFn == nil {
		return &entity.Permission{ID: uuid.New(), Code: code, Description: "Test"}, nil
	}
	return m.findByCodeForUpdFn(ctx, code)
}
func (m *mockPermRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(ctx, id)
}
func (m *mockPermRepo) FindByID(_ context.Context, _ uuid.UUID) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (m *mockPermRepo) FindByIDs(_ context.Context, _ []uuid.UUID) ([]entity.Permission, error) {
	return nil, nil
}
func (m *mockPermRepo) FindByCode(_ context.Context, _ string) (*entity.Permission, error) {
	return nil, domain.ErrPermissionNotFound
}
func (m *mockPermRepo) FindAll(_ context.Context) ([]entity.Permission, error) {
	return []entity.Permission{}, nil
}

type mockRolePermRepo struct {
	existsByPermIDFn func(ctx context.Context, permID uuid.UUID) (bool, error)
}

func (m *mockRolePermRepo) ExistsByPermissionID(ctx context.Context, permID uuid.UUID) (bool, error) {
	if m.existsByPermIDFn == nil {
		return false, nil
	}
	return m.existsByPermIDFn(ctx, permID)
}
func (m *mockRolePermRepo) FindByRoleID(_ context.Context, _ uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) FindByRoleIDs(_ context.Context, _ []uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (m *mockRolePermRepo) Assign(_ context.Context, _ entity.RolePermission) error { return nil }
func (m *mockRolePermRepo) Exists(_ context.Context, _, _ uuid.UUID) (bool, error)  { return false, nil }
func (m *mockRolePermRepo) ExistsByRoleID(_ context.Context, _ uuid.UUID) (bool, error) {
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

type permBuilder struct {
	permRepo     *mockPermRepo
	rolePermRepo *mockRolePermRepo
}

func newPermBuilder() *permBuilder {
	return &permBuilder{
		permRepo:     &mockPermRepo{},
		rolePermRepo: &mockRolePermRepo{},
	}
}

func (b *permBuilder) build() *permusecase.PermissionService {
	return permusecase.NewPermissionService(b.permRepo, b.rolePermRepo, &mockUUID{}, &inlineTx{})
}

// ─── CreatePermission ─────────────────────────────────────────────────────────

func TestCreatePermission_Success(t *testing.T) {
	b := newPermBuilder()
	var created bool
	b.permRepo.createFn = func(_ context.Context, _ entity.Permission) error {
		created = true
		return nil
	}

	svc := b.build()
	view, err := svc.CreatePermission(context.Background(), input.CreatePermissionInput{
		Code: "users.read", Description: "Read users",
	})

	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, "users.read", view.Code)
	assert.Equal(t, "Read users", view.Description)
}

func TestCreatePermission_InvalidCode(t *testing.T) {
	svc := newPermBuilder().build()
	_, err := svc.CreatePermission(context.Background(), input.CreatePermissionInput{
		Code: "", Description: "Read users",
	})
	assert.ErrorIs(t, err, domain.ErrInvalidPermissionCode)
}

func TestCreatePermission_EmptyDescription(t *testing.T) {
	svc := newPermBuilder().build()
	_, err := svc.CreatePermission(context.Background(), input.CreatePermissionInput{
		Code: "users.read", Description: "",
	})
	assert.ErrorIs(t, err, domain.ErrPermissionDescriptionEmpty)
}

func TestCreatePermission_DuplicateCode(t *testing.T) {
	b := newPermBuilder()
	b.permRepo.existsByCodeFn = func(_ context.Context, _ string) (bool, error) {
		return true, nil
	}
	svc := b.build()
	_, err := svc.CreatePermission(context.Background(), input.CreatePermissionInput{
		Code: "users.read", Description: "Read users",
	})
	assert.ErrorIs(t, err, domain.ErrDuplicatePermissionCode)
}

func TestCreatePermission_CodeNormalized(t *testing.T) {
	b := newPermBuilder()
	var createdCode string
	b.permRepo.createFn = func(_ context.Context, p entity.Permission) error {
		createdCode = p.Code
		return nil
	}

	svc := b.build()
	_, err := svc.CreatePermission(context.Background(), input.CreatePermissionInput{
		Code: "USERS.READ", Description: "Read users",
	})
	require.NoError(t, err)
	assert.Equal(t, "users.read", createdCode) // нормализован к lowercase
}

// ─── DeletePermission ─────────────────────────────────────────────────────────

func TestDeletePermission_Success(t *testing.T) {
	b := newPermBuilder()
	var deleted bool
	b.permRepo.deleteFn = func(_ context.Context, _ uuid.UUID) error {
		deleted = true
		return nil
	}

	svc := b.build()
	err := svc.DeletePermission(context.Background(), "users.read")
	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeletePermission_InvalidCode(t *testing.T) {
	svc := newPermBuilder().build()
	err := svc.DeletePermission(context.Background(), "")
	assert.ErrorIs(t, err, domain.ErrInvalidPermissionCode)
}

func TestDeletePermission_NotFound(t *testing.T) {
	b := newPermBuilder()
	b.permRepo.findByCodeForUpdFn = func(_ context.Context, _ string) (*entity.Permission, error) {
		return nil, domain.ErrPermissionNotFound
	}
	svc := b.build()
	err := svc.DeletePermission(context.Background(), "ghost.perm")
	assert.Error(t, err)
}

func TestDeletePermission_InUse_Blocked(t *testing.T) {
	b := newPermBuilder()
	b.rolePermRepo.existsByPermIDFn = func(_ context.Context, _ uuid.UUID) (bool, error) {
		return true, nil
	}
	svc := b.build()
	err := svc.DeletePermission(context.Background(), "users.read")
	assert.ErrorIs(t, err, domain.ErrPermissionInUse)
}
