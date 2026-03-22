package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, token entity.PasswordResetToken) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := q.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}

	return nil
}

func (r *PasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	return r.findByTokenHash(ctx, tokenHash, false)
}

func (r *PasswordResetRepository) FindByTokenHashForUpdate(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	return r.findByTokenHash(ctx, tokenHash, true)
}

func (r *PasswordResetRepository) findByTokenHash(ctx context.Context, tokenHash string, forUpdate bool) (*entity.PasswordResetToken, error) {
	q := ExtractTx(ctx, r.db)

	query := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1`

	if forUpdate {
		query += ` FOR UPDATE`
	}

	var t entity.PasswordResetToken
	err := q.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResetTokenNotFound
		}
		return nil, fmt.Errorf("find password reset token by hash: %w", err)
	}

	return &t, nil
}

func (r *PasswordResetRepository) InvalidateByUserID(ctx context.Context, userID uuid.UUID, now time.Time) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		UPDATE password_reset_tokens
		SET used_at = $2
		WHERE user_id = $1
		  AND used_at IS NULL`

	_, err := q.ExecContext(ctx, query, userID, now)
	if err != nil {
		return fmt.Errorf("invalidate password reset tokens by user id: %w", err)
	}

	return nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		UPDATE password_reset_tokens
		SET used_at = $2
		WHERE id = $1
		  AND used_at IS NULL`

	res, err := q.ExecContext(ctx, query, id, usedAt)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark password reset token used rows affected: %w", err)
	}
	// 0 rows означает что token уже был помечен used параллельным запросом
	if n == 0 {
		return domain.ErrResetTokenUsed
	}

	return nil
}
