package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL   = "https://chess-results.com"
	defaultUserAgent = "Mozilla/5.0 (compatible; MapleBot/1.0; +https://maple.chess)"
	defaultDelay     = 2 * time.Second
	defaultTimeout   = 15 * time.Second
	maxRetries       = 3
)

// Client is a thread-safe HTTP client for chess-results.com with rate limiting.
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	delay      time.Duration
	mu         sync.Mutex
	lastReq    time.Time
	logger     *slog.Logger
}

// Option configures the Client.
type Option func(*Client)

func WithBaseURL(url string) Option    { return func(c *Client) { c.baseURL = url } }
func WithDelay(d time.Duration) Option { return func(c *Client) { c.delay = d } }
func WithUserAgent(ua string) Option   { return func(c *Client) { c.userAgent = ua } }
func WithLogger(l *slog.Logger) Option { return func(c *Client) { c.logger = l } }

// New creates a new Client.
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		userAgent:  defaultUserAgent,
		delay:      defaultDelay,
		logger:     slog.Default(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get fetches a page, respecting rate limits and retrying on failure.
func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	c.rateLimit()

	url := c.baseURL + path
	c.logger.Debug("fetching", "url", url)

	var lastErr error
	for attempt := range maxRetries {
		body, err := c.doGet(ctx, url)
		if err == nil {
			return body, nil
		}
		lastErr = err
		c.logger.Warn("request failed, retrying", "attempt", attempt+1, "error", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.delay * time.Duration(attempt+1)):
		}
	}
	return nil, fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) doGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return body, nil
}

func (c *Client) rateLimit() {
	c.mu.Lock()
	defer c.mu.Unlock()
	elapsed := time.Since(c.lastReq)
	if elapsed < c.delay {
		time.Sleep(c.delay - elapsed)
	}
	c.lastReq = time.Now()
}

// URL builders

func TournamentURL(tournamentID string, art int) string {
	return fmt.Sprintf("/tnr%s.aspx?lan=1&art=%d", tournamentID, art)
}

func TournamentExcelURL(tournamentID string, art int) string {
	return fmt.Sprintf("/tnr%s.aspx?lan=1&art=%d&excel=2002", tournamentID, art)
}

func FederationURL(fed string) string {
	return fmt.Sprintf("/fed.aspx?lan=1&fed=%s", fed)
}
