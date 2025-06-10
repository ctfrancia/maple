package rest

import "github.com/go-chi/chi/v5"

type Router struct {
}

func NewRouter() *chi.Mux {
	routes := &Router{}

	return routes.Routes()
}

func (r *Router) Routes() *chi.Mux {
	mux := chi.NewMux()

	return mux
}
