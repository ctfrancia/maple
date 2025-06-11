package rest

import (
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	sysHealthHandler ports.SystemHealthHandler
}

func NewRouter(shh ports.SystemHealthServicer) *chi.Mux {
	routes := &Router{
		sysHealthHandler: handlers.NewSystemHealthHandler(shh),
	}

	return routes.Routes()
}

func (r *Router) Routes() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(router chi.Router) {
		mux.Get("/systemhealth", r.sysHealthHandler.Handle)
	})

	return mux
}
