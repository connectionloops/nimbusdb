package health

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// HealthPath is the path for the health endpoint
	HealthPath = "/health"
	// ReadinessPath is the path for the readiness endpoint
	ReadinessPath = "/ready"
)

var (
	// isReady indicates if the application is ready to serve traffic
	// Using atomic operations for thread-safe access
	isReady int32
)

// SetReady sets the readiness status of the application.
// This function is thread-safe.
func SetReady(ready bool) {
	if ready {
		atomic.StoreInt32(&isReady, 1)
	} else {
		atomic.StoreInt32(&isReady, 0)
	}
}

// IsReady returns the current readiness status.
// This function is thread-safe.
func IsReady() bool {
	return atomic.LoadInt32(&isReady) == 1
}

// StartHealthServer starts a lightweight HTTP server for health and readiness checks.
// The server runs in a separate goroutine and listens on the specified port.
//
// params:
//   - ctx: Context for graceful shutdown
//   - port: Port number to listen on
func StartHealthServer(ctx context.Context, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc(HealthPath, handleHealth)
	mux.HandleFunc(ReadinessPath, handleReadiness)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Int("port", port).Msg("Starting health check server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Health check server failed")
		}
	}()

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting down health check server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Error shutting down health check server")
		}
	}()
}

// handleHealth handles the /health endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error().Err(err).Msg("Failed to write health response")
	}
}

// handleReadiness handles the /ready endpoint
func handleReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if IsReady() {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Ready")); err != nil {
			log.Error().Err(err).Msg("Failed to write readiness response")
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte("Not Ready")); err != nil {
			log.Error().Err(err).Msg("Failed to write readiness response")
		}
	}
}
