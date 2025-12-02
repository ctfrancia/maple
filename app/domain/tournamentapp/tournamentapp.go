package tournamentapp

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
)

type app struct {
	tournamentBus tournamentbus.ExtBusiness
}

func newApp(tBus tournamentbus.ExtBusiness) *app {
	return &app{
		tournamentBus: tBus,
	}
}

func (a *app) create(w http.ResponseWriter, r *http.Request) {
	var nt NewTournament
	if err := json.NewDecoder(r.Body).Decode(&nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := toBusNewTournament(r.Context(), nt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tnmt, err := a.tournamentBus.Create(r.Context(), t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toAppTournament(tnmt))
}

/*
func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	return nil
}
*/
