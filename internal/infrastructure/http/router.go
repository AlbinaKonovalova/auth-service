package httpinfra

import (
	"net/http"

	httpmiddleware "github.com/AlbinaKonovalova/auth-service/internal/infrastructure/middleware/http"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

func BuildRouter(
	gwMux http.Handler,
	tokens output.TokenProvider,
	openapiJSON []byte,
) http.Handler {
	root := http.NewServeMux()

	authMW := httpmiddleware.Auth(tokens)

	protectedGw := authMW(gwMux)

	// protected: /api/v1/auth/me — must be registered before /api/v1/auth/
	root.Handle("/api/v1/auth/me", protectedGw)

	// public: all other /api/v1/auth/* and /healthz
	root.Handle("/api/v1/auth/", gwMux)
	root.Handle("/healthz", gwMux)

	// admin: all /api/v1/admin/* — protected
	root.Handle("/api/v1/admin/", protectedGw)

	if openapiJSON != nil {
		root.Handle("/swagger/", http.StripPrefix("/swagger", swaggerHandler(openapiJSON)))
	}

	return root
}
