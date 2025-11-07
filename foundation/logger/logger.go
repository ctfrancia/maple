package logger

import "log/slog"

type Logger struct {
	discard   bool
	handler   slog.Handler
	traceIDFn TraceIDFn
}
