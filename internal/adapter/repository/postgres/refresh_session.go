package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type RefreshSessionRepository struct {
	db *sql.DB
}

func NewRefreshSessionRepository(db *sql.DB) *RefreshSessionRepository {
	return &RefreshSessionRepository{db: db}
}

func (r *RefreshSessionRepository) Save(ctx context.Context, s entity.RefreshSession) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		INSERT INTO refresh_sessions (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := q.ExecContext(ctx, query,
		s.ID, s.UserID, s.TokenHash.String(), s.ExpiresAt, s.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save refresh session: %w", err)
	}

	return nil
}

func (r *RefreshSessionRepository) FindByTokenHash(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error) {
	return r.findByTokenHash(ctx, hash, false)
}

func (r *RefreshSessionRepository) FindByTokenHashForUpdate(ctx context.Context, hash value.TokenHash) (*entity.RefreshSession, error) {
	return r.findByTokenHash(ctx, hash, true)
}

func (r *RefreshSessionRepository) findByTokenHash(ctx context.Context, hash value.TokenHash, forUpdate bool) (*entity.RefreshSession, error) {
	q := ExtractTx(ctx, r.db)

	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_sessions
		WHERE token_hash = $1`

	if forUpdate {
		query += ` FOR UPDATE`
	}

	var s entity.RefreshSession
	var tokenHash string
	err := q.QueryRowContext(ctx, query, hash.String()).Scan(
		&s.ID, &s.UserID, &tokenHash, &s.ExpiresAt, &s.RevokedAt, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("find refresh session by token hash: %w", err)
	}

	s.TokenHash = value.TokenHash(tokenHash)

	return &s, nil
}

func (r *RefreshSessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		UPDATE refresh_sessions
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL`

	_, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("revoke refresh session: %w", err)
	}

	return nil
}

func (r *RefreshSessionRepository) DeleteExpiredAndRevoked(ctx context.Context, userID uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		DELETE FROM refresh_sessions
		WHERE user_id = $1
		  AND (expires_at <= now() OR revoked_at IS NOT NULL)`

	_, err := q.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete expired/revoked sessions: %w", err)
	}

	return nil
}

func (r *RefreshSessionRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	q := ExtractTx(ctx, r.db)

	const query = `
		UPDATE refresh_sessions
		SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()`

	_, err := q.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke all sessions by user id: %w", err)
	}

	return nil
}
