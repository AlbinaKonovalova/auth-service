package authservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	authadapter "github.com/AlbinaKonovalova/auth-service/internal/adapter/auth"
	"github.com/AlbinaKonovalova/auth-service/internal/adapter/controller"
	pgadapter "github.com/AlbinaKonovalova/auth-service/internal/adapter/repository/postgres"
	"github.com/AlbinaKonovalova/auth-service/internal/adapter/system"
	"github.com/AlbinaKonovalova/auth-service/internal/config"
	httpinfra "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/http"
	pginfra "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/postgres"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/auth"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/user"
	authservicepkg "github.com/AlbinaKonovalova/auth-service/pkg/authservice"
)

type App struct {
	server *httpinfra.Server
	db     *sql.DB
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := pginfra.NewConnection(pginfra.Config{
		URL:             cfg.Database.URL,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	userRepo := pgadapter.NewUserRepository(db)
	roleRepo := pgadapter.NewRoleRepository(db)
	permRepo := pgadapter.NewPermissionRepository(db)
	userRoleRepo := pgadapter.NewUserRoleRepository(db)
	rolePermRepo := pgadapter.NewRolePermissionRepository(db)
	sessionRepo := pgadapter.NewRefreshSessionRepository(db)
	txManager := pgadapter.NewTxManager(db)

	hasher := authadapter.NewPasswordHasher(authadapter.Argon2Params{
		Memory:      cfg.Argon2.Memory,
		Iterations:  cfg.Argon2.Iterations,
		Parallelism: cfg.Argon2.Parallelism,
		SaltLength:  cfg.Argon2.SaltLength,
		KeyLength:   cfg.Argon2.KeyLength,
	})
	tokenProvider := authadapter.NewTokenProvider(authadapter.TokenProviderConfig{
		Secret:            cfg.Auth.JWTSecret,
		AccessTokenTTL:    cfg.Auth.AccessTokenTTL,
		RefreshTokenBytes: cfg.Auth.RefreshTokenBytes,
	})
	tokenHasher := authadapter.NewTokenHasher()
	clockImpl := system.NewClock()
	uuidGen := system.NewUUIDGenerator()

	resolver := common.NewPermissionResolver(userRoleRepo, rolePermRepo, roleRepo, permRepo)

	authService := auth.NewAuthService(userRepo, sessionRepo, hasher, tokenProvider, tokenHasher, txManager, clockImpl, uuidGen, resolver, cfg.Cookie.TTL)

	userService := user.NewUserService(userRepo, roleRepo, userRoleRepo, sessionRepo, hasher, txManager, clockImpl, uuidGen)

	ctrl := controller.NewAuthServiceController(authService, userService, cfg.Cookie, logger)

	ctx := context.Background()
	gwMux, err := httpinfra.NewGatewayMux(ctx, ctrl)
	if err != nil {
		return nil, fmt.Errorf("build gateway: %w", err)
	}

	router := httpinfra.BuildRouter(gwMux, tokenProvider, authservicepkg.SwaggerSpec)

	server := httpinfra.NewServer(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		logger,
		router,
		cfg.CORS.AllowedOrigins,
	)

	return &App{server: server, db: db, logger: logger}, nil
}

func (a *App) Run() error {
	if err := a.server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	a.logger.Info("closing postgres connection")
	return a.db.Close()
}
