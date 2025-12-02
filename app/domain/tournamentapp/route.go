package tournamentapp

import (
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log           *logger.Logger
	TournamentBus tournamentbus.ExtBusiness
	//AuthClient    *auth.Client
}

func Routes(cfg Config) chi.Router {
	const v1 = "v1"
	api := newApp(cfg.TournamentBus)
	r := chi.NewRouter()

	// Domain-specific middleware
	//r.Use(middleware.TournamentAuth)

	r.Route("/"+v1, func(v1 chi.Router) {
		v1.Post("/tournament", api.create)
		//v1.Get("/tournaments", api.query)
		//v1.Get("/tournament/{id}", api.queryByID)
	})

	return r
}
