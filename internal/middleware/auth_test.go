package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Mock handler that will be wrapped by the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("success")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	tests := []struct {
		name           string
		apiKeyEnv      string
		requestAPIKey  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid API key",
			apiKeyEnv:      "test-api-key",
			requestAPIKey:  "test-api-key",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid API key",
			apiKeyEnv:      "test-api-key",
			requestAPIKey:  "wrong-api-key",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid API key",
		},
		{
			name:           "Missing API key in request",
			apiKeyEnv:      "test-api-key",
			requestAPIKey:  "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Missing API key",
		},
		{
			name:           "API key not configured",
			apiKeyEnv:      "",
			requestAPIKey:  "some-key",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "API key not configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.apiKeyEnv != "" {
				if err := os.Setenv("API_KEY", tt.apiKeyEnv); err != nil {
					t.Fatalf("failed to set API_KEY: %v", err)
				}
				defer func() {
					if err := os.Unsetenv("API_KEY"); err != nil {
						t.Errorf("failed to unset API_KEY: %v", err)
					}
				}()
			} else {
				if err := os.Unsetenv("API_KEY"); err != nil {
					t.Fatalf("failed to unset API_KEY: %v", err)
				}
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.requestAPIKey != "" {
				req.Header.Set("X-API-Key", tt.requestAPIKey)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Apply middleware and serve
			handler := AuthMiddleware(mockHandler)
			handler.ServeHTTP(w, req)

			// Check status code
			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Check error response if expected
			if tt.expectedError != "" {
				var errResp ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}

				if errResp.Error != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, errResp.Error)
				}

				contentType := resp.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %q", contentType)
				}
			}
		})
	}
}

func TestWriteUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Test error message"

	writeUnauthorized(w, message)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error != message {
		t.Errorf("expected error %q, got %q", message, errResp.Error)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}
