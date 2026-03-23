package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	authadapter "github.com/AlbinaKonovalova/auth-service/internal/adapter/auth"
	"github.com/AlbinaKonovalova/auth-service/internal/config"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	email := flag.String("email", "admin@example.com", "admin email (default: admin@example.com)")
	password := flag.String("password", "", "admin password (required)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if *password == "" {
		logger.Error("password is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		logger.Error("failed to open db", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping db", slog.Any("error", err))
		os.Exit(1)
	}

	hasher := authadapter.NewPasswordHasher(cfg.Argon2)

	ctx := context.Background()

	if err := seedAdmin(ctx, db, hasher, *email, *password, logger); err != nil {
		logger.Error("seed failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("seed complete")
}

func seedAdmin(
	ctx context.Context,
	db *sql.DB,
	hasher *authadapter.PasswordHasher,
	rawEmail, rawPassword string,
	logger *slog.Logger,
) error {
	email, err := value.NewEmail(rawEmail)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if _, err := value.NewPassword(rawPassword); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	hash, err := hasher.Hash(rawPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var userID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, is_active)
		VALUES ($1, $2, true)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id`,
		email.String(), hash,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}

	var roleID string
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM roles WHERE code = 'admin'`,
	).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("role 'admin' not found: run migrations first")
		}
		return fmt.Errorf("get admin role: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, roleID,
	)
	if err != nil {
		return fmt.Errorf("assign admin role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	logger.Info("admin user seeded",
		slog.String("email", email.String()),
		slog.String("user_id", userID),
	)

	return nil
}
