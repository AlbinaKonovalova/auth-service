package httpinfra

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed swagger-ui
var swaggerUI embed.FS

func swaggerHandler(openapiJSON []byte) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openapiJSON) //nolint:errcheck
	})

	uiFS, _ := fs.Sub(swaggerUI, "swagger-ui")
	fileServer := http.FileServer(http.FS(uiFS))
	mux.Handle("GET /", fileServer)

	return mux
}
