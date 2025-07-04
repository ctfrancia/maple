// rest package is the REST API implementation
package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ctfrancia/maple/internal/adapters/rest/handlers"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	sysHandler ports.SystemHandler
}

func NewRouter(ss ports.SystemServicer, logger ports.LoggerServicer) *chi.Mux {
	routes := &Router{
		sysHandler: handlers.NewSystemHandler(ss, logger),
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
	})

	// TODO: should only print if not in production
	printRoutes(mux)

	return mux
}

func printRoutes(r chi.Router) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%-6s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
