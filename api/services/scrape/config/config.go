package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ardanlabs/conf/v3"
)

type Tempo struct {
	Host        string  `conf:"default:tempo:4317"`
	ServiceName string  `conf:"default:maple"`
	Probability float64 `conf:"default:0.5"`
	// Shouldn't use a high Probability value in non-developer systems.
	// 0.05 should be enough for most systems. Some might want to have
	// this even lower.
}

type Web struct {
	ReadTimeout        time.Duration `conf:"default:5s"`
	WriteTimeout       time.Duration `conf:"default:10s"`
	IdleTimeout        time.Duration `conf:"default:120s"`
	ShutdownTimeout    time.Duration `conf:"default:20s"`
	APIHost            string        `conf:"default:0.0.0.0:8080"`
	DebugHost          string        `conf:"default:0.0.0.0:3010"`
	CORSAllowedOrigins []string      `conf:"default:*"`
}

// Config holds all service configuration, loaded from environment variables.
type Config struct {
	// Service
	conf.Version
	Web         Web
	ServiceName string `conf:"default:chess-results-scraper"`
	Environment string `conf:"default:development"`
	LogLevel    string `conf:"default:info"`

	// HTTP server
	//HTTPPort int `conf:"default:8080"`

	// Chess-results client
	BaseURL   string        `conf:"default:https://chess-results.com"`
	UserAgent string        `conf:"default:Mozilla/5.0 (compatible; MapleBot/1.0; +https://maple.chess)"`
	ReqDelay  time.Duration `conf:"default:2s"`

	// Scheduler
	DiscoveryCron     string        `conf:"default:6h"` // Cron expression for federation discovery (e.g. "0 */6 * * *")
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
	/*
		cfg := struct {
			conf.Version
			Web struct {
				ReadTimeout        time.Duration `conf:"default:5s"`
				WriteTimeout       time.Duration `conf:"default:10s"`
				IdleTimeout        time.Duration `conf:"default:120s"`
				ShutdownTimeout    time.Duration `conf:"default:20s"`
				APIHost            string        `conf:"default:0.0.0.0:3000"`
				DebugHost          string        `conf:"default:0.0.0.0:3010"`
				CORSAllowedOrigins []string      `conf:"default:*"`
			}
			Auth struct {
				Host string `conf:"default:http://auth-service:6000"`
			}
			DB struct {
				User         string `conf:"default:postgres"`
				Password     string `conf:"default:postgres,mask"`
				Host         string `conf:"default:database-service"`
				Name         string `conf:"default:postgres"`
				MaxIdleConns int    `conf:"default:0"`
				MaxOpenConns int    `conf:"default:0"`
				DisableTLS   bool   `conf:"default:true"`
			}
			Tempo struct {
				Host        string  `conf:"default:tempo:4317"`
				ServiceName string  `conf:"default:maple"`
				Probability float64 `conf:"default:0.5"`
				// Shouldn't use a high Probability value in non-developer systems.
				// 0.05 should be enough for most systems. Some might want to have
				// this even lower.
			}
		}{
			Version: conf.Version{
				Build: build,
				Desc:  "Maple",
			},
		}
	*/

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

	const prefix = "MAPLE"
	help, err := conf.Parse(prefix, c)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, nil
		}

		return nil, fmt.Errorf("parsing config: %w", err)
	}

	/*
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
	*/

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
