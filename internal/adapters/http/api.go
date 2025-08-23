// Package http is the REST API implementation
package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ctfrancia/maple/internal/adapters/http/handlers/system"
	"github.com/ctfrancia/maple/internal/adapters/http/handlers/tournament"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	sysHandler        ports.SystemHandler
	tournamentHandler ports.TournamentHandler
}

func NewRouter(log ports.Logger, ss ports.SystemServicer, ts ports.TournamentServicer) *chi.Mux {
	routes := &Router{
		sysHandler:        systemhandlers.NewSystemHandler(ss, log),
		tournamentHandler: tournamenthandlers.NewTournamentHandler(log, ts),
	}

	return routes.Routes()
}

func (r *Router) Routes() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(v1 chi.Router) {
		v1.Route("/system", func(v1s chi.Router) {
			v1s.Get("/health", r.sysHandler.HealthHandler)
			v1s.Post("/login", r.sysHandler.LoginHandler)
			v1s.Post("/new-consumer", r.sysHandler.NewConsumerHandler)
		})
		v1.Route("/tournament", func(v1t chi.Router) {
			v1t.Get("/find/{id}", r.tournamentHandler.FindTournamentHandler)
			v1t.Post("/new", r.tournamentHandler.CreateTournamentHandler)
			// v1t.Post("/tournaments", r.tournamentHandler.CreateTournamentHandler)
			// v1t.Put("/tournaments/{id}", r.tournamentHandler.UpdateTournamentHandler)
			// v1t.Delete("/tournaments/{id}", r.tournamentHandler.DeleteTournamentHandler)
		})
		v1.Route("/match", func(v1m chi.Router) {
			// v1m.Get("/matches", r.matchHandler.GetMatchesHandler)
			// v1m.Post("/matches", r.tournamentHandler.CreateMatchHandler)
			// v1m.Put("/matches/{id}", r.matchHandler.UpdateMatchHandler)
			// v1m.Delete("/matches/{id}", r.matchHandler.DeleteMatchHandler)
		})
	})

	// TODO: should only print if not in production
	printRoutes(mux)

	return mux
}

func printRoutes(r chi.Router) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		fmt.Printf("%-6s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
