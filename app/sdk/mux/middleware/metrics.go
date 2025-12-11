package middleware

import (
	//"github.com/ctfrancia/maple/foundation/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"})
)

/*
func Otel(tracer trace.Tracer) func(next http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.InjectTracing(r.Context(), tracer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return h // return the handler
	}
	return m // return the middleware
}
*/

// Metrics middleware records metrics with exemplars attached.
func Metrics() func(next http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// process request
			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(ww.statusCode)

			// Extract trace ID from span context for exemplar
			span := trace.SpanFromContext(r.Context())
			traceID := span.SpanContext().TraceID().String()

			fmt.Printf("DEBUG Metrics: traceID=%s, valid=%v\n", traceID, span.SpanContext().IsValid())

			labels := prometheus.Labels{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": status,
			}

			// Record metrics with exemplar linking to trace
			if traceID != "" && span.SpanContext().IsValid() {
				exemplarLabels := prometheus.Labels{"traceID": traceID}
				observer := httpRequestDuration.With(labels)

				if exemplarObs, ok := observer.(prometheus.ExemplarObserver); ok {
					exemplarObs.ObserveWithExemplar(duration, exemplarLabels)
					fmt.Printf("DEBUG: Recorded exemplar successfully\n")
				} else {
					fmt.Printf("DEBUG: Failed to cast to ExemplarObserver\n")
					observer.Observe(duration)
				}
			} else {
				fmt.Printf("DEBUG: No valid traceID, using regular observe\n")
				httpRequestDuration.With(labels).Observe(duration)
			}

			httpRequestsTotal.With(labels).Inc()
		})
		return h // return the handler
	}
	return m
}

/*
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// process request
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.statusCode)

		// Extract trace ID from span context for exemplar
		span := trace.SpanFromContext(r.Context())
		traceID := span.SpanContext().TraceID().String()

		labels := prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": status,
		}

		// Record metrics with exemplar linking to trace
		if traceID != "" {
			exemplarLabels := prometheus.Labels{"traceID": traceID}
			httpRequestDuration.With(labels).(prometheus.ExemplarObserver).
				ObserveWithExemplar(duration, exemplarLabels)
		} else {
			httpRequestDuration.With(labels).Observe(duration)
		}

		httpRequestsTotal.With(labels).Inc()
	})
}
*/
