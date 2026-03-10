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
