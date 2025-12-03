// Package mux provides suport to bind domain to http handlers.
package mux

import (
	"net/http"

	"github.com/ctfrancia/maple/app/sdk/mux/middleware"
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/web"

	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type BusConfig struct {
	TournamentBus tournamentbus.ExtBusiness
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
	app := web.NewChiApp(
		cfg.Log.Info,
		cfg.Tracer,
	)
	// Apply general middleware to the Chi router
	app.Use(middleware.Otel(cfg.Tracer)) // need to create span first
	app.Use(middleware.Metrics)          // metrics has exemplars
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.CORSAllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	app.Use(chimid.Recoverer)

	return app
}
