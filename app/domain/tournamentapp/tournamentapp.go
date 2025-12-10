package tournamentapp

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/app/sdk/web"
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

// create creates a new tournament.
// func (a *app) create(w http.ResponseWriter, r *http.Request) {
func (a *app) create(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	var nt NewTournament
	if err := json.NewDecoder(r.Body).Decode(&nt); err != nil {
		return &web.Response{Status: http.StatusBadRequest}, err
	}

	t, err := toBusNewTournament(r.Context(), nt)
	if err != nil {
		return &web.Response{Status: http.StatusBadRequest}, err
	}

	tnmt, err := a.tournamentBus.Create(r.Context(), t)
	if err != nil {
		return nil, err
	}

	return &web.Response{
		Data:   toAppTournament(tnmt),
		Status: http.StatusOK,
	}, nil
}

// query retrieves a list of tournaments based on the provided filter.
// TODO: implement query
func (a *app) query(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	date := r.URL.Query().Get("date")
	player := r.URL.Query().Get("player")
	if date != "" && player != "" {
		return nil, nil
	}

	return nil, nil
}

// fetch retrieves a tournament by ID.
func (a *app) fetch(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	return nil, nil
}

// delete removes a tournament by ID.
func (a *app) delete(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	return nil, nil
}

// update updates a tournament by ID.
func (a *app) update(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	return nil, nil
}

// partialUpdate updates a tournament by ID with the provided values.
func (a *app) partialUpdate(w http.ResponseWriter, r *http.Request) (*web.Response, error) {
	return nil, nil
}
