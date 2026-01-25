package tournamentapp

import (
	"github.com/ctfrancia/maple/app/sdk/web"
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log           *logger.Logger
	TournamentBus tournamentbus.ExtBusiness
	//AuthClient    *auth.Client
}

func V1Routes(cfg Config) chi.Router {
	r := chi.NewRouter()
	api := newApp(cfg.TournamentBus)
	v1BasePath := "/v1/tournament"

	// Domain-specific middleware
	//r.Use(middleware.TournamentAuth)

	r.Route(v1BasePath, func(v1 chi.Router) {
		v1.Post("/", web.Wrap(cfg.Log, api.create))
		v1.Get("/", web.Wrap(cfg.Log, api.query))
		v1.Get("/{uuid}", web.Wrap(cfg.Log, api.fetch))
		v1.Put("/{uuid}", web.Wrap(cfg.Log, api.update))
		v1.Patch("/{uuid}", web.Wrap(cfg.Log, api.partialUpdate))
		v1.Delete("/{uuid}", web.Wrap(cfg.Log, api.delete))
	})

	return r
}
