package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

func (r *RolePermissionRepository) FindByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.RolePermission, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT role_id, permission_id
		FROM role_permissions
		WHERE role_id = $1`

	rows, err := q.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("find role permissions by role id: %w", err)
	}
	defer rows.Close()

	var result []entity.RolePermission
	for rows.Next() {
		var rp entity.RolePermission
		if err := rows.Scan(&rp.RoleID, &rp.PermissionID); err != nil {
			return nil, fmt.Errorf("scan role permission: %w", err)
		}
		result = append(result, rp)
	}

	return result, rows.Err()
}

func (r *RolePermissionRepository) FindByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]entity.RolePermission, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	q := ExtractTx(ctx, r.db)

	uuidStrings := make([]string, len(roleIDs))
	for i, id := range roleIDs {
		uuidStrings[i] = id.String()
	}

	const query = `
		SELECT role_id, permission_id
		FROM role_permissions
		WHERE role_id = ANY($1)`

	rows, err := q.QueryContext(ctx, query, pq.Array(uuidStrings))
	if err != nil {
		return nil, fmt.Errorf("find role permissions by role ids: %w", err)
	}
	defer rows.Close()

	var result []entity.RolePermission
	for rows.Next() {
		var rp entity.RolePermission
		if err := rows.Scan(&rp.RoleID, &rp.PermissionID); err != nil {
			return nil, fmt.Errorf("scan role permission: %w", err)
		}
		result = append(result, rp)
	}

	return result, rows.Err()
}

func (r *RolePermissionRepository) Assign(ctx context.Context, rp entity.RolePermission) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	_, err := q.ExecContext(ctx, query, rp.RoleID, rp.PermissionID)
	if err != nil {
		return fmt.Errorf("assign permission to role: %w", err)
	}

	return nil
}

func (r *RolePermissionRepository) Exists(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT 1 FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2
		LIMIT 1`

	var dummy int
	err := q.QueryRowContext(ctx, query, roleID, permissionID).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check role permission exists: %w", err)
	}

	return true, nil
}

func (r *RolePermissionRepository) Revoke(ctx context.Context, roleID, permissionID uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2`

	_, err := q.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("revoke permission from role: %w", err)
	}

	return nil
}

func (r *RolePermissionRepository) ExistsByRoleID(ctx context.Context, roleID uuid.UUID) (bool, error) {
	q := ExtractTx(ctx, r.db)

	const query = `SELECT EXISTS(SELECT 1 FROM role_permissions WHERE role_id = $1)`

	var exists bool
	if err := q.QueryRowContext(ctx, query, roleID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check role permission exists by role id: %w", err)
	}

	return exists, nil
}

func (r *RolePermissionRepository) ExistsByPermissionID(ctx context.Context, permissionID uuid.UUID) (bool, error) {
	q := ExtractTx(ctx, r.db)

	const query = `SELECT EXISTS(SELECT 1 FROM role_permissions WHERE permission_id = $1)`

	var exists bool
	if err := q.QueryRowContext(ctx, query, permissionID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check role permission exists by permission id: %w", err)
	}

	return exists, nil
}
