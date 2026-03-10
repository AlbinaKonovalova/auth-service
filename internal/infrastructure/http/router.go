package httpinfra

import (
	"net/http"

	httpmiddleware "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/middleware/http"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

// BuildRouter собирает root HTTP mux с явным разделением public/protected routes.
//
// Public routes (без auth):
//
//	POST /auth/login
//	POST /auth/refresh
//	POST /auth/logout
//	GET  /healthz
//	GET  /swagger/...
//
// Protected routes (требуют Bearer access token):
//
//	GET /auth/me
func BuildRouter(
	gwMux http.Handler,
	tokens output.TokenProvider,
	openapiJSON []byte,
) http.Handler {
	root := http.NewServeMux()

	authMW := httpmiddleware.Auth(tokens)

	// protectedGw — gateway с auth middleware.
	// Используется только для защищённых маршрутов.
	protectedGw := authMW(gwMux)

	// Protected routes — явно перечислены, не перекрываются через /auth/.
	root.Handle("GET /auth/me", protectedGw)

	// Public gateway routes.
	root.Handle("/auth/", gwMux)
	root.Handle("/healthz", gwMux)

	// Swagger UI + spec.
	if openapiJSON != nil {
		root.Handle("/swagger/", http.StripPrefix("/swagger", swaggerHandler(openapiJSON)))
	}

	return root
}
