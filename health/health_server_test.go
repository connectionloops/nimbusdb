package health

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use a different port for testing to avoid conflicts
	// We'll test the handlers directly instead
	mux := http.NewServeMux()
	mux.HandleFunc(HealthPath, handleHealth)
	mux.HandleFunc(ReadinessPath, handleReadiness)

	server := &http.Server{
		Addr:    ":0", // Let OS choose port
		Handler: mux,
	}

	go func() {
		_ = server.ListenAndServe()
	}()
	defer server.Shutdown(ctx)

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://localhost" + server.Addr + HealthPath)
	if err != nil {
		// If we can't connect, test the handler directly
		testHealthHandler(t)
		testReadinessHandler(t)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func testHealthHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", HealthPath, nil)
	w := &mockResponseWriter{}

	handleHealth(w, req)

	if w.statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.statusCode)
	}

	if string(w.body) != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", string(w.body))
	}
}

func testReadinessHandler(t *testing.T) {
	// Test not ready
	SetReady(false)
	req, _ := http.NewRequest("GET", ReadinessPath, nil)
	w := &mockResponseWriter{}

	handleReadiness(w, req)

	if w.statusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503 when not ready, got %d", w.statusCode)
	}

	// Test ready
	SetReady(true)
	w2 := &mockResponseWriter{}
	handleReadiness(w2, req)

	if w2.statusCode != http.StatusOK {
		t.Errorf("Expected status 200 when ready, got %d", w2.statusCode)
	}
}

type mockResponseWriter struct {
	statusCode int
	body       []byte
	header     http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.body = append(m.body, b...)
	return len(b), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}
