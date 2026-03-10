package httpinfra

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed swagger-ui
var swaggerUI embed.FS

// swaggerHandler returns an http.Handler that serves Swagger UI.
func swaggerHandler(openapiJSON []byte) http.Handler {
	mux := http.NewServeMux()

	// Serve openapi.json
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openapiJSON) //nolint:errcheck
	})

	// Serve swagger-ui static files
	uiFS, _ := fs.Sub(swaggerUI, "swagger-ui")
	fileServer := http.FileServer(http.FS(uiFS))
	mux.Handle("GET /", fileServer)

	return mux
}
