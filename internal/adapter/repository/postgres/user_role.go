package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

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
