// Package debug provides handler support for the debugging endpoints.
package debug

import (
	"expvar"
	//"net/http"
	//"net/http/pprof"

	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	// "github.com/gin-gonic/gin"
)

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
/*
func Mux() *http.ServeMux {
	r := chi.NewRouter()
	srv, _ := statsviz.NewServer()

	r.Get("/debug/statsviz/ws", srv.Ws())
	//r.HandleFunc("/debug/pprof/", pprof.Index)
	//r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	//r.Handle("/debug/vars/", expvar.Handler())
	r.Handle("/debug/statsviz/*", srv.Index())

	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars/", expvar.Handler())

	//statsviz.Register(m)

	return r
}
*/

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
// FIXME: this is not working
func Mux() *chi.Mux {
	r := chi.NewRouter()
	srv, _ := statsviz.NewServer() // we aren't passing any opts so won't error

	r.Get("/debug/statsviz/ws", srv.Ws())
	// r.HandleFunc("/debug/pprof/", pprof.Index)
	// r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// r.Handle("/debug/vars/", expvar.Handler())
	r.Handle("/debug/statsviz/*", srv.Index())
	r.Handle("/debug/vars", expvar.Handler())

	return r
}
