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

	// Web server
	APIHost            string
	DebugHost          string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	CORSAllowedOrigins []string

	// Chess-results client
	BaseURL   string
	UserAgent string
	ReqDelay  time.Duration

	// Scheduler
	DiscoveryCron     string
	DetailCron        string
	DiscoveryInterval time.Duration
	DetailInterval    time.Duration

	// Federations to monitor
	Federations []string

	// Location filter (optional — filter tournaments by location substring)
	LocationFilter string

	// Store
	DatabaseURL string

	// Publisher
	PublisherType string
	WebhookURL    string
	NatsURL       string
	NatsSubject   string
}

// Load reads config from environment variables with sensible defaults.
func Load() (*Config, error) {
	c := &Config{
		ServiceName:        envOr("SCRAPER_SERVICE_NAME", "chess-results-scraper"),
		Environment:        envOr("SCRAPER_ENVIRONMENT", "development"),
		LogLevel:           envOr("SCRAPER_LOG_LEVEL", "info"),
		APIHost:            envOr("SCRAPER_API_HOST", "0.0.0.0:8080"),
		DebugHost:          envOr("SCRAPER_DEBUG_HOST", "0.0.0.0:8081"),
		ReadTimeout:        envDurationOr("SCRAPER_READ_TIMEOUT", 5*time.Second),
		WriteTimeout:       envDurationOr("SCRAPER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:        envDurationOr("SCRAPER_IDLE_TIMEOUT", 120*time.Second),
		ShutdownTimeout:    envDurationOr("SCRAPER_SHUTDOWN_TIMEOUT", 20*time.Second),
		CORSAllowedOrigins: envListOr("SCRAPER_CORS_ALLOWED_ORIGINS", []string{"*"}),
		BaseURL:            envOr("SCRAPER_CHESS_RESULTS_BASE_URL", "https://chess-results.com"),
		UserAgent:          envOr("SCRAPER_USER_AGENT", "Mozilla/5.0 (compatible; MapleBot/1.0; +https://maple.chess)"),
		ReqDelay:           envDurationOr("SCRAPER_REQUEST_DELAY", 2*time.Second),
		DiscoveryCron:      envOr("SCRAPER_DISCOVERY_CRON", ""),
		DetailCron:         envOr("SCRAPER_DETAIL_CRON", ""),
		DiscoveryInterval:  envDurationOr("SCRAPER_DISCOVERY_INTERVAL", 6*time.Hour),
		DetailInterval:     envDurationOr("SCRAPER_DETAIL_INTERVAL", 1*time.Hour),
		Federations:        envListOr("SCRAPER_FEDERATIONS", []string{"ESP"}),
		LocationFilter:     envOr("SCRAPER_LOCATION_FILTER", ""),
		DatabaseURL:        envOr("DATABASE_URL", ""),
		PublisherType:      envOr("SCRAPER_PUBLISHER_TYPE", "log"),
		WebhookURL:         envOr("SCRAPER_WEBHOOK_URL", ""),
		NatsURL:            envOr("SCRAPER_NATS_URL", ""),
		NatsSubject:        envOr("SCRAPER_NATS_SUBJECT", "chess.tournaments"),
	}

	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if len(c.Federations) == 0 {
		return fmt.Errorf("SCRAPER_FEDERATIONS must have at least one entry")
	}
	if c.PublisherType == "webhook" && c.WebhookURL == "" {
		return fmt.Errorf("SCRAPER_WEBHOOK_URL required when SCRAPER_PUBLISHER_TYPE=webhook")
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
