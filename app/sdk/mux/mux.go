// Package mux provides suport to bind domain to http handlers.
package mux

import (
	"net/http"

	"github.com/ctfrancia/maple/app/domain/scraperapp"
	"github.com/ctfrancia/maple/app/domain/tournamentapp"
	"github.com/ctfrancia/maple/app/sdk/mux/middleware"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type BusConfig struct {
	TournamentBus tournamentbus.ExtBusiness
	ScraperBus    scraperbus.ExtBusiness
}

// AuthConfig contains all the mandarory components for the system.
// TODO: implement auth
type AuthConfig struct {
	//Auth auth.Auth
}

// Config contains all the mandarory components for the system.
type Config struct {
	CORSAllowedOrigins []string
	Build              string
	Tracer             trace.Tracer
	Log                *logger.Logger
	DB                 *gorm.DB
	BusConfig          BusConfig
	AuthConfig         AuthConfig
}

func WebAPI(cfg Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Otel(cfg.Tracer))
	r.Use(middleware.LoggingMiddleware(cfg.Log))
	r.Use(middleware.Metrics())
	r.Use(chimid.Recoverer)

	// Mount tournament routes if configured
	if cfg.BusConfig.TournamentBus != nil {
		tCfg := tournamentapp.Config{Log: cfg.Log, TournamentBus: cfg.BusConfig.TournamentBus}
		r.Mount("/api", tournamentapp.V1Routes(tCfg))
	}

	// Mount scraper routes if configured
	if cfg.BusConfig.ScraperBus != nil {
		sCfg := scraperapp.Config{Log: cfg.Log, ScraperBus: cfg.BusConfig.ScraperBus}
		r.Mount("/api", scraperapp.V1Routes(sCfg))
	}

	return r
}
