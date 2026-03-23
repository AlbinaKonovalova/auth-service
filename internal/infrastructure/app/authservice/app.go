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
	mailadapter "github.com/AlbinaKonovalova/auth-service/internal/adapter/repository/mail"
	pgadapter "github.com/AlbinaKonovalova/auth-service/internal/adapter/repository/postgres"
	"github.com/AlbinaKonovalova/auth-service/internal/adapter/system"
	"github.com/AlbinaKonovalova/auth-service/internal/config"
	httpinfra "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/http"
	pginfra "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/postgres"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/access"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/auth"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/common"
	passwordreset "github.com/AlbinaKonovalova/auth-service/internal/usecase/password_reset"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/permission"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/role"
	"github.com/AlbinaKonovalova/auth-service/internal/usecase/user"
	authservicepkg "github.com/AlbinaKonovalova/auth-service/pkg/authservice"
)

type App struct {
	server *httpinfra.Server
	db     *sql.DB
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := pginfra.NewConnection(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	userRepo := pgadapter.NewUserRepository(db)
	roleRepo := pgadapter.NewRoleRepository(db)
	permRepo := pgadapter.NewPermissionRepository(db)
	userRoleRepo := pgadapter.NewUserRoleRepository(db)
	rolePermRepo := pgadapter.NewRolePermissionRepository(db)
	sessionRepo := pgadapter.NewRefreshSessionRepository(db)
	resetRepo := pgadapter.NewPasswordResetRepository(db)
	txManager := pgadapter.NewTxManager(db)

	hasher := authadapter.NewPasswordHasher(cfg.Argon2)
	tokenProvider := authadapter.NewTokenProvider(cfg.Auth)
	tokenHasher := authadapter.NewTokenHasher()
	clockImpl := system.NewClock()
	uuidGen := system.NewUUIDGenerator()

	mailer := mailadapter.NewSMTPMailer(cfg.SMTP, cfg.PasswordReset)

	resolver := common.NewPermissionResolver(userRoleRepo, rolePermRepo, roleRepo, permRepo)

	authService := auth.NewAuthService(userRepo, sessionRepo, hasher, tokenProvider, tokenHasher, txManager, clockImpl, uuidGen, resolver, cfg.Cookie.TTL)

	userService := user.NewUserService(userRepo, roleRepo, userRoleRepo, sessionRepo, hasher, txManager, clockImpl, uuidGen)

	accessService := access.NewAccessService(userRepo, userRoleRepo, roleRepo, rolePermRepo, permRepo, clockImpl, txManager)

	roleService := role.NewRoleService(roleRepo, userRoleRepo, rolePermRepo, uuidGen, txManager)

	permissionService := permission.NewPermissionService(permRepo, rolePermRepo, uuidGen, txManager)

	passwordResetService := passwordreset.NewPasswordResetService(
		userRepo, resetRepo, sessionRepo, hasher, tokenProvider, tokenHasher, mailer, clockImpl, uuidGen, txManager, cfg.PasswordReset.TTL,
	)

	ctrl := controller.NewAuthServiceController(
		authService, userService, accessService, roleService, permissionService, passwordResetService, cfg.Cookie, logger,
	)

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
