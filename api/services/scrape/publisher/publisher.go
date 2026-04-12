package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ctfrancia/maple/api/services/scrape/events"
)

// Publisher emits events when tournaments are discovered or updated.
// Other maple microservices can consume these events.
type Publisher interface {
	Publish(ctx context.Context, event events.Event) error
	Close() error
}

// ---------- Log publisher (development) ----------

// LogPublisher logs events to slog — useful for dev and debugging.
type LogPublisher struct {
	logger *slog.Logger
}

func NewLogPublisher(logger *slog.Logger) *LogPublisher {
	return &LogPublisher{logger: logger}
}

func (p *LogPublisher) Publish(_ context.Context, event events.Event) error {
	p.logger.Info("event published",
		"type", event.Type,
		"tournament_id", event.TournamentID,
		"federation", event.Federation,
		"timestamp", event.Timestamp,
	)
	return nil
}

func (p *LogPublisher) Close() error { return nil }

// ---------- Webhook publisher ----------

// WebhookPublisher sends events as JSON POST requests to a webhook URL.
// This is a simple way to integrate with other services (e.g., maple API).
type WebhookPublisher struct {
	url    string
	client *http.Client
	logger *slog.Logger
}

func NewWebhookPublisher(url string, logger *slog.Logger) *WebhookPublisher {
	return &WebhookPublisher{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

func (p *WebhookPublisher) Publish(ctx context.Context, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-Type", string(event.Type))

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	p.logger.Debug("webhook delivered", "type", event.Type, "tournament_id", event.TournamentID)
	return nil
}

func (p *WebhookPublisher) Close() error { return nil }

// ---------- Multi publisher ----------

// MultiPublisher fans out events to multiple publishers.
type MultiPublisher struct {
	publishers []Publisher
}

func NewMultiPublisher(publishers ...Publisher) *MultiPublisher {
	return &MultiPublisher{publishers: publishers}
}

func (p *MultiPublisher) Publish(ctx context.Context, event events.Event) error {
	var firstErr error
	for _, pub := range p.publishers {
		if err := pub.Publish(ctx, event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (p *MultiPublisher) Close() error {
	for _, pub := range p.publishers {
		pub.Close()
	}
	return nil
}
