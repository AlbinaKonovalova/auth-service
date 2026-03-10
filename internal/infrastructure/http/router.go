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

	root.Handle("GET /auth/me", protectedGw)

	root.Handle("/auth/", gwMux)
	root.Handle("/healthz", gwMux)

	if openapiJSON != nil {
		root.Handle("/swagger/", http.StripPrefix("/swagger", swaggerHandler(openapiJSON)))
	}

	return root
}
