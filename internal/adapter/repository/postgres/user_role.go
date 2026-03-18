package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type UserRoleRepository struct {
	db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserRole, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT user_id, role_id, created_at
		FROM user_roles
		WHERE user_id = $1`

	rows, err := q.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find user roles: %w", err)
	}
	defer rows.Close()

	var result []entity.UserRole
	for rows.Next() {
		var ur entity.UserRole
		if err := rows.Scan(&ur.UserID, &ur.RoleID, &ur.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		result = append(result, ur)
	}

	return result, rows.Err()
}

func (r *UserRoleRepository) FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]entity.UserRole, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	q := ExtractTx(ctx, r.db)

	strs := make([]string, len(userIDs))
	for i, id := range userIDs {
		strs[i] = id.String()
	}

	const query = `
		SELECT user_id, role_id, created_at
		FROM user_roles
		WHERE user_id = ANY($1)`

	rows, err := q.QueryContext(ctx, query, pq.Array(strs))
	if err != nil {
		return nil, fmt.Errorf("find user roles by user ids: %w", err)
	}
	defer rows.Close()

	var result []entity.UserRole
	for rows.Next() {
		var ur entity.UserRole
		if err := rows.Scan(&ur.UserID, &ur.RoleID, &ur.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		result = append(result, ur)
	}

	return result, rows.Err()
}

func (r *UserRoleRepository) Assign(ctx context.Context, ur entity.UserRole) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO user_roles (user_id, role_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`

	_, err := q.ExecContext(ctx, query, ur.UserID, ur.RoleID, ur.CreatedAt)
	if err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}
	return nil
}
