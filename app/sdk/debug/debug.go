// Package debug provides handler support for the debugging endpoints.
package debug

import (
	"expvar"
	"net/http/pprof"

	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
func Mux() *chi.Mux {
	r := chi.NewRouter()
	srv, _ := statsviz.NewServer() // we aren't passing any opts so won't error

	r.Get("/debug/statsviz/ws", srv.Ws())
	r.Get("/debug/pprof/", pprof.Index)
	r.Get("/debug/pprof/cmdline", pprof.Cmdline)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/symbol", pprof.Symbol)
	r.Get("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/statsviz/*", srv.Index())
	r.Handle("/debug/vars", expvar.Handler())
	//r.Handle("/metrics", promhttp.Handler())
	// Enable exemplars in the handler
	r.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true, // This enables exemplar support
		},
	))
	return r
}
