package tournamentapp

import (
	"net/http"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/web"
)

type Config struct {
	Log           *logger.Logger
	TournamentBus tournamentbus.Business
}

func Routes(app *web.App, cfg Config) {
	const v1 = "v1"

	// TODO: create middlewares for authentication and authorization
	// mid.authentication()
	// mid.authorization()

	api := newApp(cfg.TournamentBus)

	app.HandlerFunc(http.MethodGet, v1, "/tournaments", api.query)
	app.HandlerFunc(http.MethodPost, v1, "/tournaments", api.create)
}
