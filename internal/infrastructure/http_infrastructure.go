// Package infrastructure contains the infrastructure for the API
// the api contains the shared use of making 3rd party API calls.
package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ctfrancia/maple/internal/core/ports"
)

// HTTPClientAdapter is an adapter for the http.Client
type HTTPClientAdapter struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClientAdapter creates a new HTTPClientAdapter
func NewHTTPClientAdapter(baseURL string, timeout time.Duration) *HTTPClientAdapter {
	return &HTTPClientAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// Get makes a GET request to the given URL
func (h *HTTPClientAdapter) Get(ctx context.Context, url string, headers map[string]string) (*ports.HTTPResponse, error) {
	return h.makeRequest(ctx, http.MethodGet, url, nil, headers)
}

// Post makes a POST request to the given URL
func (h *HTTPClientAdapter) Post(ctx context.Context, url string, body []byte, headers map[string]string) (*ports.HTTPResponse, error) {
	return h.makeRequest(ctx, http.MethodPost, url, body, headers)
}

// Put makes a PUT request to the given URL
func (h *HTTPClientAdapter) Put(ctx context.Context, url string, body []byte, headers map[string]string) (*ports.HTTPResponse, error) {
	return h.makeRequest(ctx, http.MethodPut, url, body, headers)
}

// Delete makes a DELETE request to the given URL
func (h *HTTPClientAdapter) Delete(ctx context.Context, url string, headers map[string]string) (*ports.HTTPResponse, error) {
	return h.makeRequest(ctx, http.MethodDelete, url, nil, headers)
}

func (h *HTTPClientAdapter) makeRequest(ctx context.Context, method, url string, body []byte, headers map[string]string) (*ports.HTTPResponse, error) {
	fullURL := h.baseURL + url

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	return &ports.HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    respHeaders,
	}, nil
}
