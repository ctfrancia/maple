package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientAdapter_Get(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		url            string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
		expectedError  bool
	}{
		{
			name: "successful GET request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/test", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "success"}`))
			},
			url: "/test",
			headers: map[string]string{
				"Accept": "application/json",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "success"}`,
			expectedError:  false,
		},
		{
			name: "GET request with 404 response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "not found"}`))
			},
			url:            "/notfound",
			headers:        map[string]string{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "not found"}`,
			expectedError:  false,
		},
		{
			name: "GET request with multiple headers",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
				assert.Equal(t, "my-app/1.0", r.Header.Get("User-Agent"))

				w.Header().Set("X-Response-ID", "12345")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data": "test"}`))
			},
			url: "/headers",
			headers: map[string]string{
				"Accept":        "application/json",
				"Authorization": "Bearer token123",
				"User-Agent":    "my-app/1.0",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data": "test"}`,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Create adapter
			adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

			// Make request
			ctx := context.Background()
			resp, err := adapter.Get(ctx, tt.url, tt.headers)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(resp.Body))
		})
	}
}

func TestHTTPClientAdapter_Post(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		url            string
		body           []byte
		headers        map[string]string
		expectedStatus int
		expectedBody   string
		expectedError  bool
	}{
		{
			name: "successful POST request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/create", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify request body
				var requestBody map[string]any
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				require.NoError(t, err)
				assert.Equal(t, "test", requestBody["name"])

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"id": "123", "name": "test"}`))
			},
			url:  "/create",
			body: []byte(`{"name": "test"}`),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id": "123", "name": "test"}`,
			expectedError:  false,
		},
		{
			name: "POST request with empty body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, int64(0), r.ContentLength)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "empty body received"}`))
			},
			url:            "/empty",
			body:           nil,
			headers:        map[string]string{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "empty body received"}`,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

			ctx := context.Background()
			resp, err := adapter.Post(ctx, tt.url, tt.body, tt.headers)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(resp.Body))
		})
	}
}

func TestHTTPClientAdapter_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/update/123", r.URL.Path)

		var requestBody map[string]any
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		assert.Equal(t, "updated", requestBody["status"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "status": "updated"}`))
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

	ctx := context.Background()
	body := []byte(`{"status": "updated"}`)
	headers := map[string]string{"Content-Type": "application/json"}

	resp, err := adapter.Put(ctx, "/update/123", body, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"id": "123", "status": "updated"}`, string(resp.Body))
}

func TestHTTPClientAdapter_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/delete/123", r.URL.Path)
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

	ctx := context.Background()
	headers := map[string]string{"Authorization": "Bearer token123"}

	resp, err := adapter.Delete(ctx, "/delete/123", headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Empty(t, resp.Body)
}

func TestHTTPClientAdapter_ContextCancellation(t *testing.T) {
	// Create a server that sleeps longer than our context timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := adapter.Get(ctx, "/slow", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestHTTPClientAdapter_Timeout(t *testing.T) {
	// Create a server that sleeps longer than client timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create adapter with short timeout
	adapter := NewHTTPClientAdapter(server.URL, 100*time.Millisecond)

	ctx := context.Background()
	_, err := adapter.Get(ctx, "/slow", nil)

	assert.Error(t, err)
	// assert.Contains(t, err.Error(), "timeout")

	assert.True(t,
		assert.ObjectsAreEqual("Client.Timeout exceeded", err.Error()) ||
			assert.ObjectsAreEqual("context deadline exceeded", err.Error()) ||
			assert.Contains(t, err.Error(), "Client.Timeout exceeded") ||
			assert.Contains(t, err.Error(), "context deadline exceeded"),
		"Expected timeout-related error, got: %v", err)
}

func TestHTTPClientAdapter_InvalidURL(t *testing.T) {
	adapter := NewHTTPClientAdapter("", 5*time.Second)

	ctx := context.Background()
	_, err := adapter.Get(ctx, "://invalid-url", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestHTTPClientAdapter_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

	ctx := context.Background()
	resp, err := adapter.Get(ctx, "/error", nil)

	require.NoError(t, err) // HTTP errors don't return Go errors
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"error": "internal server error"}`, string(resp.Body))
}

func TestHTTPClientAdapter_ResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-Id", "req-123")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)

	ctx := context.Background()
	resp, err := adapter.Get(ctx, "/headers", nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that we have headers
	assert.NotEmpty(t, resp.Headers)

	// Check specific headers (case-sensitive as they come from the server)
	assert.Equal(t, "application/json", resp.Headers["Content-Type"])
	assert.Equal(t, "req-123", resp.Headers["X-Request-Id"])
	assert.Equal(t, "no-cache", resp.Headers["Cache-Control"])

	// Debug output to help diagnose issues
	if resp.Headers["X-Request-Id"] == "" {
		t.Logf("Available headers: %+v", resp.Headers)
		for k, v := range resp.Headers {
			t.Logf("Header: '%s' = '%s'", k, v)
		}
	}
}

func TestHTTPClientAdapter_LargeResponse(t *testing.T) {
	// Create a large response body
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(largeData)
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 10*time.Second)

	ctx := context.Background()
	resp, err := adapter.Get(ctx, "/large", nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, len(largeData), len(resp.Body))
	assert.Equal(t, largeData, resp.Body)
}

// Benchmark tests
func BenchmarkHTTPClientAdapter_Get(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "benchmark"}`))
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := adapter.Get(ctx, "/benchmark", nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTTPClientAdapter_Post(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123"}`))
	}))
	defer server.Close()

	adapter := NewHTTPClientAdapter(server.URL, 5*time.Second)
	ctx := context.Background()
	body := []byte(`{"name": "benchmark"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := adapter.Post(ctx, "/benchmark", body, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Integration test example
func TestHTTPClientAdapter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This would typically hit a real API endpoint
	// For demo purposes, we'll use httpbin.org (if available)
	adapter := NewHTTPClientAdapter("https://httpbin.org", 10*time.Second)

	ctx := context.Background()

	t.Run("real GET request", func(t *testing.T) {
		resp, err := adapter.Get(ctx, "/get", map[string]string{
			"User-Agent": "go-test-client",
		})

		if err != nil {
			t.Skipf("Skipping integration test due to network error: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, string(resp.Body), "httpbin")
	})

	t.Run("real POST request", func(t *testing.T) {
		body := []byte(`{"test": "data"}`)
		resp, err := adapter.Post(ctx, "/post", body, map[string]string{
			"Content-Type": "application/json",
		})

		if err != nil {
			t.Skipf("Skipping integration test due to network error: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, string(resp.Body), "test")
	})
}
