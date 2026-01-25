// Package middleware provides middleware for the mux.
package middleware

import (
	"github.com/ctfrancia/maple/foundation/otel"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func Otel(tracer trace.Tracer) func(next http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), "http.request")
			defer span.End()

			ctx = otel.InjectTracing(ctx, tracer)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return h // return the handler
	}
	return m // return the middleware
}
