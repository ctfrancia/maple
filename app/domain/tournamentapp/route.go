package tournamentapp

import (
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/web"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    any // TODO: define DB interface
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp()
}
