package httpinfra

import (
	"net/http"

	middleware "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/middleware/http"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

func BuildRouter(
	gwMux http.Handler,
	tokens output.TokenProvider,
	openapiJSON []byte,
) http.Handler {
	root := http.NewServeMux()

	authMW := middleware.Auth(tokens)

	protectedGw := authMW(gwMux)
	adminProtectedGw := authMW(middleware.AdminPermissionGuard(gwMux))

	// protected: /api/v1/auth/me — must be registered before /api/v1/auth/
	root.Handle("/api/v1/auth/me", protectedGw)

	// public: password-reset routes — registered before /api/v1/auth/ for clarity
	// /api/v1/auth/password-reset/request and /api/v1/auth/password-reset/confirm
	// are public: no auth middleware, no permission guard
	root.Handle("/api/v1/auth/password-reset/", gwMux)

	// public: all other /api/v1/auth/* and /healthz
	root.Handle("/api/v1/auth/", gwMux)
	root.Handle("/healthz", gwMux)

	// admin: all /api/v1/admin/* — adminProtectedGw
	root.Handle("/api/v1/admin/", adminProtectedGw)

	if openapiJSON != nil {
		root.Handle("/swagger/", http.StripPrefix("/swagger", swaggerHandler(openapiJSON)))
	}

	return root
}
