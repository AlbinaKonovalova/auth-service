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
	email := flag.String("email", "", "admin email (required)")
	password := flag.String("password", "", "admin password (required)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if *email == "" || *password == "" {
		logger.Error("email and password are required")
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

	hasher := authadapter.NewPasswordHasher(authadapter.Argon2Params{
		Memory:      cfg.Argon2.Memory,
		Iterations:  cfg.Argon2.Iterations,
		Parallelism: cfg.Argon2.Parallelism,
		SaltLength:  cfg.Argon2.SaltLength,
		KeyLength:   cfg.Argon2.KeyLength,
	})

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
	// Валидируем email через value object — нормализует и проверяет формат.
	email, err := value.NewEmail(rawEmail)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	// Валидируем пароль через value object — применяет доменные правила.
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

	// Идемпотентный INSERT: если пользователь уже существует — пропускаем,
	// без предварительного SELECT (исключаем гонку).
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

	// Получаем role_id для admin.
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

	// Назначаем роль — ON CONFLICT делает это идемпотентным.
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
