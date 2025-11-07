package web

import (
	"context"
	"net/http"
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
	tracer  any
	mux     *http.ServeMux
	otmux   http.Handler
	mw      []MidFunc
	origins []string
}
