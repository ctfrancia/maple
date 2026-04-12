package scraperapp

import (
	"context"

	"github.com/ctfrancia/maple/app/sdk/web"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Orchestrator defines the interface for triggering scrapes
type Orchestrator interface {
	TriggerScrape(ctx context.Context, federation string) (uuid.UUID, error)
}

type Config struct {
	Log          *logger.Logger
	ScraperBus   scraperbus.ExtBusiness
	Orchestrator Orchestrator
}

func V1Routes(cfg Config) chi.Router {
	r := chi.NewRouter()
	api := newApp(cfg.ScraperBus, cfg.Orchestrator)
	v1BasePath := "/v1/scraper"

	r.Route(v1BasePath, func(v1 chi.Router) {
		// Scrape job endpoints
		v1.Post("/jobs", web.Wrap(cfg.Log, api.createJob))
		v1.Get("/jobs", web.Wrap(cfg.Log, api.queryJobs))
		v1.Get("/jobs/{id}", web.Wrap(cfg.Log, api.fetchJob))

		// Scraped tournament endpoints
		v1.Get("/tournaments", web.Wrap(cfg.Log, api.queryTournaments))

		// Manual trigger
		v1.Post("/trigger", web.Wrap(cfg.Log, api.triggerScrape))

		// Statistics
		v1.Get("/statistics", web.Wrap(cfg.Log, api.getStatistics))
	})

	return r
}
