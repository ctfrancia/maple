package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

// TraceIDFn represents a functionthat can return the trace id from
// the specified context.
type TraceIDFn func(ctx context.Context) string

type Logger struct {
	discard   bool
	handler   slog.Handler
	traceIDFn TraceIDFn
}

// New constructs a new log for application use.
func New(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, Events{})
}

// NewWithEvents contructs a new log for application use with events
func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, events)
}

// NewWithHandler returns a new log for application use with the underlying
// handler.
func NewWithHandler(h slog.Handler) *Logger {
	return &Logger{handler: h}
}

// NewStdLogger returns a standard library logger that wraps the slog logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(logger.handler, slog.Level(level))
}

// Debug logs at LevelDebug with the given context.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelDebug, 3, msg, args...)
}

// Debugc logs the information at the specified call stack postion. (c = caller)
func (l *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelDebug, caller, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack postion. (c = caller)
func (l *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelInfo, caller, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelWarn, 3, msg, args...)
}

// Warnc logs the information at the specified call stack postion. (c = caller)
func (l *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelWarn, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelError, 3, msg, args...)
}

// Errorc logs the information at the specified call stack postion. (c = caller)
func (l *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	if l.discard {
		return
	}

	l.write(ctx, LevelError, caller, msg, args...)
}

func new(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.Level(minLevel), ReplaceAttr: f}))

	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		handler = newLogHandler(handler, events)
	}

	// Attributes to add to all logs.
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	// Add those attrubutes and capture the final handler.
	handler = handler.WithAttrs(attrs)

	return &Logger{
		discard:   w == io.Discard,
		handler:   handler,
		traceIDFn: traceIDFn,
	}
}

func (l *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	slogLevel := slog.Level(level)

	if !l.handler.Enabled(ctx, slogLevel) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), slogLevel, msg, pcs[0])

	if l.traceIDFn != nil {
		args = append(args, "trace_id", l.traceIDFn(ctx))
	}

	r.Add(args...)

	l.handler.Handle(ctx, r)

}
