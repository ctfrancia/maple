package tournamentapp

import (
	"context"
	"net/http"

	"github.com/ctfrancia/maple/foundation/web"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
)

type app struct {
	//tournamentBus tournamentbus.TournamentBus
}

func newApp(tBus tournamentbus.Business) *app {
	return &app{
		//tournamentBus: tBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}
