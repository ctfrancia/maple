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

func V1Routes(cfg Config) chi.Router {
	r := chi.NewRouter()
	api := newApp(cfg.TournamentBus)

	// Domain-specific middleware
	//r.Use(middleware.TournamentAuth)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/tournament", api.create)
	})

	return r
}
