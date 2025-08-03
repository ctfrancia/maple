package ports

import "context"

type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, url string, body []byte, headers map[string]string) (*HTTPResponse, error)
	Put(ctx context.Context, url string, body []byte, headers map[string]string) (*HTTPResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}
