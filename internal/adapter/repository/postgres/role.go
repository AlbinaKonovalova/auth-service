package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
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

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*entity.Role, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, code, name, description
		FROM roles
		WHERE code = $1`

	var role entity.Role
	err := q.QueryRowContext(ctx, query, code).Scan(&role.ID, &role.Code, &role.Name, &role.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, fmt.Errorf("find role by code: %w", err)
	}

	return &role, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, code, name, description
		FROM roles`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
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

func (r *RoleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	q := ExtractTx(ctx, r.db)

	const query = `SELECT EXISTS(SELECT 1 FROM roles WHERE code = $1)`

	var exists bool
	if err := q.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		return false, fmt.Errorf("check role exists by code: %w", err)
	}

	return exists, nil
}

func (r *RoleRepository) Create(ctx context.Context, role entity.Role) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO roles (id, code, name, description)
		VALUES ($1, $2, $3, $4)`

	if _, err := q.ExecContext(ctx, query, role.ID, role.Code, role.Name, role.Description); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrDuplicateRoleCode
		}
		return fmt.Errorf("create role: %w", err)
	}

	return nil
}
