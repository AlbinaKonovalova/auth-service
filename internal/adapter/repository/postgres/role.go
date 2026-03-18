package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Role, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := ExtractTx(ctx, r.db)

	uuidStrings := make([]string, len(ids))
	for i, id := range ids {
		uuidStrings[i] = id.String()
	}

	const query = `
		SELECT id, code, name, description
		FROM roles
		WHERE id = ANY($1)`

	rows, err := q.QueryContext(ctx, query, pq.Array(uuidStrings))
	if err != nil {
		return nil, fmt.Errorf("find roles by ids: %w", err)
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Code, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *RoleRepository) FindByCodes(ctx context.Context, codes []string) ([]entity.Role, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, code, name, description
		FROM roles
		WHERE code = ANY($1)`

	rows, err := q.QueryContext(ctx, query, pq.Array(codes))
	if err != nil {
		return nil, fmt.Errorf("find roles by codes: %w", err)
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Code, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("scan role by code: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}
