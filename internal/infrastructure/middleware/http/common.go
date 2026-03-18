package middleware

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

type apiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func apiError(code int, message string) apiErrorResponse {
	return apiErrorResponse{
		Code:    code,
		Message: message,
	}
}
