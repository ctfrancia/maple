package tournamentapp

import (
	//"net/http"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"
	//"github.com/ctfrancia/maple/foundation/web"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log           *logger.Logger
	TournamentBus tournamentbus.Business
}

/*
func Routes(app *web.App, cfg Config) {
	const v1 = "v1"

	// TODO: create middlewares for authentication and authorization
	// mid.authentication()
	// mid.authorization()

	api := newApp(cfg.TournamentBus)

	app.HandlerFunc(http.MethodGet, v1, "/tournaments", api.query)
	app.HandlerFunc(http.MethodPost, v1, "/tournaments", api.create)
}
*/

func Routes(cfg Config) chi.Router {
	const v1 = "v1"

	api := newApp(cfg.TournamentBus)

	r := chi.NewRouter()
	r.Route(v1, func(v1 chi.Router) {
		//v1.Get("/tournaments", api.query)
		v1.Post("/tournaments", api.create)
	})

	return r
}
