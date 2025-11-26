package prometheus

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ctfrancia/maple/foundation/logger"
)

// Exporter implements the prometheus exporter interface
type Exporter struct {
	log    *logger.Logger
	server http.Server
	data   map[string]any
	mu     sync.Mutex
}

func New(log *logger.Logger, host, route string, readTO, writeTO, idleTO time.Duration) *Exporter {
	mux := http.NewServeMux()

	exp := Exporter{
		log: log,
		server: http.Server{
			Addr:         host,
			Handler:      mux,
			ReadTimeout:  readTO,
			WriteTimeout: writeTO,
			IdleTimeout:  idleTO,
			ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
		},
	}

	mux.HandleFunc(route, exp.handler)

	go func() {
		ctx := context.Background()
		log.Info(ctx, "prometheus", "satus", "API Listening", "host", host)

		if err := exp.server.ListenAndServe(); err != nil {
			log.Error(ctx, "prometheus", "err", err)
		}
	}()

	return &exp
}

func (e *Exporter) Publish(data map[string]any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.data = deepCopyMap(data)
}

func (e *Exporter) Stop(shutdownTO time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTO)
	defer cancel()

	e.log.Info(ctx, "prometheus", "status", "start shutdown...")
}

func (e *Exporter) handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	var data map[string]any
	e.mu.Lock()
	{
		data = deepCopyMap(e.data)
	}
	e.mu.Unlock()

	out(w, "", data)

	e.log.Info(ctx, "prometheus", "metrics", fmt.Sprintf("expvar : (%d) : %s %s -> %s", http.StatusOK, r.Method, r.URL.Path, r.RemoteAddr))
}

func deepCopyMap(source map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range source {
		switch vm := v.(type) {
		case map[string]any:
			result[k] = deepCopyMap(vm)

		case int64:
			result[k] = float64(vm)

		case float64:
			result[k] = vm

		case bool:
			result[k] = 0.0
			if vm {
				result[k] = 1.0
			}
		}
	}

	return result
}

func out(w io.Writer, prefix string, data map[string]any) {
	if prefix != "" {
		prefix += "_"
	}

	for k, v := range data {
		writeKey := fmt.Sprintf("%s%s", prefix, k)

		switch vm := v.(type) {
		case float64:
			fmt.Fprintf(w, "%s %.f\n", writeKey, vm)

		case map[string]any:
			out(w, writeKey, vm)

		default:
			// Discard this value.
		}
	}
}
