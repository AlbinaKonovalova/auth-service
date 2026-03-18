package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/output"
)

func Auth(tokens output.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := bearerFromRequest(r)
			if raw == "" {
				writeJSON(w, http.StatusUnauthorized, apiError(16, "missing authorization header"))
				return
			}

			claims, err := tokens.ParseAccessToken(r.Context(), raw)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, apiError(16, "invalid or expired access token"))
				return
			}

			ctx := context.WithValue(r.Context(), value.ClaimsContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerFromRequest(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if strings.HasPrefix(v, "Bearer ") {
		return v[len("Bearer "):]
	}
	return ""
}
