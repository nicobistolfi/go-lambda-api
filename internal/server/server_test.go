package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	addr := ":8080"
	srv := New(addr)

	if srv == nil {
		t.Fatal("expected server instance, got nil")
	}

	if srv.httpServer == nil {
		t.Fatal("expected http server instance, got nil")
	}

	if srv.httpServer.Addr != addr {
		t.Errorf("expected address %q, got %q", addr, srv.httpServer.Addr)
	}

	if srv.httpServer.ReadTimeout != 15*time.Second {
		t.Errorf("expected read timeout 15s, got %v", srv.httpServer.ReadTimeout)
	}

	if srv.httpServer.WriteTimeout != 15*time.Second {
		t.Errorf("expected write timeout 15s, got %v", srv.httpServer.WriteTimeout)
	}

	if srv.httpServer.IdleTimeout != 60*time.Second {
		t.Errorf("expected idle timeout 60s, got %v", srv.httpServer.IdleTimeout)
	}
}

func TestHealthEndpoint(t *testing.T) {
	srv := New(":0") // Use port 0 to let the system assign a free port

	// Test the health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.httpServer.Handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("test")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	// Wrap with logging middleware
	handler := loggingMiddleware(testHandler)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(w, req)

	// Verify response
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	// Test default status code
	if rw.statusCode != http.StatusOK {
		t.Errorf("expected default status %d, got %d", http.StatusOK, rw.statusCode)
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	if rw.statusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rw.statusCode)
	}

	// Verify it was passed to the underlying ResponseWriter
	if w.Code != http.StatusNotFound {
		t.Errorf("expected underlying writer status %d, got %d", http.StatusNotFound, w.Code)
	}
}
