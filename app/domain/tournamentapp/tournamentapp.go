package tournamentapp

import (
	"context"
	"net/http"

	"github.com/ctfrancia/maple/foundation/web"
)

type app struct {
	//tournamentBus tournamentbus.TournamentBus
}

func newApp(tBus tournamentbus.TournamentBus) *app {
	return &app{
		//tournamentBus: tBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}
