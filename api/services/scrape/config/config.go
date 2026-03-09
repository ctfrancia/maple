package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all service configuration, loaded from environment variables.
type Config struct {
	// Service
	ServiceName string
	Environment string
	LogLevel    string

	// HTTP server
	HTTPPort int

	// Chess-results client
	BaseURL   string
	UserAgent string
	ReqDelay  time.Duration

	// Scheduler
	DiscoveryCron     string        // Cron expression for federation discovery (e.g. "0 */6 * * *")
	DetailCron        string        // Cron expression for tournament detail scraping
	DiscoveryInterval time.Duration // Alternative: fixed interval for discovery
	DetailInterval    time.Duration // Alternative: fixed interval for detail scraping

	// Federations to monitor
	Federations []string // e.g. ["ESP", "CAT"] for Spain / Catalonia

	// Location filter (optional — filter tournaments by location substring)
	LocationFilter string // e.g. "Barcelona"

	// Store
	DatabaseURL string // postgres connection string

	// Publisher
	PublisherType string // "log", "webhook", "nats" (extensible)
	WebhookURL    string // if PublisherType == "webhook"
	NatsURL       string // if PublisherType == "nats"
	NatsSubject   string
}

// Load reads config from environment variables with sensible defaults.
func Load() (*Config, error) {
	c := &Config{
		ServiceName:       envOr("SERVICE_NAME", "chess-results-scraper"),
		Environment:       envOr("ENVIRONMENT", "development"),
		LogLevel:          envOr("LOG_LEVEL", "info"),
		HTTPPort:          envIntOr("HTTP_PORT", 8080),
		BaseURL:           envOr("CHESS_RESULTS_BASE_URL", "https://chess-results.com"),
		UserAgent:         envOr("CHESS_RESULTS_USER_AGENT", "Mozilla/5.0 (compatible; MapleBot/1.0; +https://maple.chess)"),
		ReqDelay:          envDurationOr("CHESS_RESULTS_REQUEST_DELAY", 2*time.Second),
		DiscoveryCron:     envOr("DISCOVERY_CRON", ""),
		DetailCron:        envOr("DETAIL_CRON", ""),
		DiscoveryInterval: envDurationOr("DISCOVERY_INTERVAL", 6*time.Hour),
		DetailInterval:    envDurationOr("DETAIL_INTERVAL", 1*time.Hour),
		Federations:       envListOr("FEDERATIONS", []string{"ESP"}),
		LocationFilter:    envOr("LOCATION_FILTER", ""),
		DatabaseURL:       envOr("DATABASE_URL", ""),
		PublisherType:     envOr("PUBLISHER_TYPE", "log"),
		WebhookURL:        envOr("WEBHOOK_URL", ""),
		NatsURL:           envOr("NATS_URL", ""),
		NatsSubject:       envOr("NATS_SUBJECT", "chess.tournaments"),
	}

	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if len(c.Federations) == 0 {
		return fmt.Errorf("FEDERATIONS must have at least one entry")
	}
	if c.PublisherType == "webhook" && c.WebhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL required when PUBLISHER_TYPE=webhook")
	}
	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envDurationOr(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func envListOr(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		var result []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return fallback
}
