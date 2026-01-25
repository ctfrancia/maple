// Package web provides a simple web framework for our application.
// most likely can deprecate this in favor of chi
package web

import (
	"context"
	"fmt"
	"net/http"

	//"github.com/ctfrancia/maple/app/sdk/mux/middleware"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Encoder defines behavior that can encode a data model and provide
// the content type for that encoding.
type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

// HandlerFunc represents a function that handles an http request within
// our little web framework.
type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

// Logger represents a function that adds info to the log.
type Logger func(ctx context.Context, msg string, args ...any)

// App is the entrypoint into our application and what configures our
// context object for each of our http handlers.
type App struct {
	log     Logger
	tracer  trace.Tracer
	mux     *http.ServeMux
	otmux   http.Handler
	mw      []MidFunc
	mids    []http.Handler
	origins []string
}

// App2 is the entrypoint into our application and what configures our
// context object for each of our http handlers.
type App2 struct {
	log     Logger
	tracer  trace.Tracer
	mux     *chi.Mux
	otmux   http.Handler
	mw      []MidFunc
	mids    []http.Handler
	origins []string
}

// NewChiApp creates an App value that handle a set of routes for the application using chi
// the reponsibility of this function is to mount the api routes
func NewChiApp(log Logger, tracer trace.Tracer) *chi.Mux {
	r := chi.NewRouter()
	const apiT = "api/tournament"

	//r.Mount(apiT, tournamentapp.Routes(tCfg))

	return r
}

func NewApp2(log Logger, tracer trace.Tracer, mw ...MidFunc) *App2 {
	mux := chi.NewRouter()

	return &App2{
		log:    log,
		tracer: tracer,
		mux:    mux,
		otmux:  otelhttp.NewHandler(mux, "request"),
		mw:     mw,
	}
}

func (a *App2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

// NewApp creates an App value that handle a set of routes for the application
// DEPRECATED - use NewChiApp
func NewApp(log Logger, tracer trace.Tracer, mw ...MidFunc) *App {
	// Create an otel http handler which wraps our router. This will
	// start the initial span and annotate it with info about the request/trusted.
	mux := http.NewServeMux()

	return &App{
		log:    log,
		tracer: tracer,
		mux:    mux,
		otmux:  otelhttp.NewHandler(mux, "request"),
		mw:     mw,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.origins != nil {
		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}

	a.otmux.ServeHTTP(w, r)
}

// EnableCors enables CORS preflight requests to work. It prevents the
// MethodNotAllowedHandler from being called
func (a *App) EnableCors(origins []string) {
	a.origins = origins
}

func (a *App) HandlerFuncNoMid(method, group, path string, handlerFunc HandlerFunc) {
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := setWriter(r.Context(), w)

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			a.log(ctx, "web-respond", err)
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}

func (a *App) HandlerFunc(method, group, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	handlerFunc = wrapMiddleware(mw, handlerFunc)
	handlerFunc = wrapMiddleware(a.mw, handlerFunc)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := setTracer(r.Context(), a.tracer)
		ctx = setWriter(ctx, w)

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			a.log(ctx, "web-respond", "ERROR", err)
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}
