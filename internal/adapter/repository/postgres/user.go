package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/lib/pq"

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

func (r *UserRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	q := ExtractTx(ctx, r.db)

	const query = `
		SELECT id, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
		FOR UPDATE`

	var u entity.User
	err := q.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id for update: %w", err)
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

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO users (id, email, password_hash, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := q.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.IsActive, user.CreatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrEmailAlreadyTaken
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, f dto.UserListFilters) ([]entity.User, int, error) {
	q := ExtractTx(ctx, r.db)

	args := []any{}
	conds := []string{}
	idx := 1

	if f.IsActive != nil {
		conds = append(conds, fmt.Sprintf("u.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	if f.Role != "" {
		conds = append(conds, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN roles r ON r.id = ur.role_id
			WHERE ur.user_id = u.id AND r.code = $%d
		)`, idx))
		args = append(args, f.Role)
		idx++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, where)
	var total int
	if err := q.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	perPage := f.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	listQuery := fmt.Sprintf(`
		SELECT u.id, u.email, u.password_hash, u.is_active, u.created_at, u.updated_at
		FROM users u
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`, where, idx, idx+1)
	args = append(args, perPage, offset)

	rows, err := q.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list users rows: %w", err)
	}

	return users, total, nil
}

func (r *UserRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash string) error {
	q := ExtractTx(ctx, r.db)

	const query = `UPDATE users SET password_hash = $2 WHERE id = $1`

	res, err := q.ExecContext(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("update password hash: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update password hash rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Activate(ctx context.Context, id uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `UPDATE users SET is_active = true WHERE id = $1`

	res, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("activate user: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("activate user rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `UPDATE users SET is_active = false WHERE id = $1`

	res, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deactivate user rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
