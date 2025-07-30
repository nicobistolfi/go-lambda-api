package middleware

import (
	"encoding/json"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			writeUnauthorized(w, "API key not configured")
			return
		}

		requestAPIKey := r.Header.Get("X-API-Key")
		if requestAPIKey == "" {
			writeUnauthorized(w, "Missing API key")
			return
		}

		if requestAPIKey != apiKey {
			writeUnauthorized(w, "Invalid API key")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		// Response already started, can't change status
		// In production, this should be logged
		_ = err
	}
}
