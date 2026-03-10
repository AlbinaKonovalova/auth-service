package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var u entity.User
	err := q.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	var u entity.User
	err := q.QueryRowContext(ctx, query, email.String()).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &u, nil
}
