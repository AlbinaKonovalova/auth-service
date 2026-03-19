package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, code, description
		FROM permissions
		WHERE id = $1`

	var p entity.Permission
	err := q.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Code, &p.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("find permission by id: %w", err)
	}

	return &p, nil
}

func (r *PermissionRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := ExtractTx(ctx, r.db)

	uuidStrings := make([]string, len(ids))
	for i, id := range ids {
		uuidStrings[i] = id.String()
	}

	const query = `SELECT id, code, description FROM permissions WHERE id = ANY($1)`

	rows, err := q.QueryContext(ctx, query, pq.Array(uuidStrings))
	if err != nil {
		return nil, fmt.Errorf("find permissions by ids: %w", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}

	return permissions, rows.Err()
}

func (r *PermissionRepository) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, code, description
		FROM permissions
		WHERE code = $1`

	var p entity.Permission
	err := q.QueryRowContext(ctx, query, code).Scan(&p.ID, &p.Code, &p.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("find permission by code: %w", err)
	}

	return &p, nil
}

func (r *PermissionRepository) FindAll(ctx context.Context) ([]entity.Permission, error) {
	q := ExtractTx(ctx, r.db)

	const query = `SELECT id, code, description FROM permissions ORDER BY code`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all permissions: %w", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}

	return permissions, rows.Err()
}

func (r *PermissionRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	q := ExtractTx(ctx, r.db)

	const query = `SELECT EXISTS(SELECT 1 FROM permissions WHERE code = $1)`

	var exists bool
	if err := q.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		return false, fmt.Errorf("check permission exists by code: %w", err)
	}

	return exists, nil
}

func (r *PermissionRepository) Create(ctx context.Context, p entity.Permission) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO permissions (id, code, description)
		VALUES ($1, $2, $3)`

	if _, err := q.ExecContext(ctx, query, p.ID, p.Code, p.Description); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrDuplicatePermissionCode
		}
		return fmt.Errorf("create permission: %w", err)
	}

	return nil
}
