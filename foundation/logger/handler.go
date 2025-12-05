package logger

import (
	"context"
	"log/slog"
)

// logHandler provides a wrapper around the slog handler to capture which log
// level is beinglogged for event handling.
type logHandler struct {
	handler slog.Handler
	events  Events
}

func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handlers records at the given level.
// The handler discards records with lower levels.
func (l *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

// WithAttrs returns a new JSONHandler whose attributes consists
// of l's attributes followed by attrs.
func (l *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: l.handler.WithAttrs(attrs), events: l.events}
}

// WithGroup returns a new Handler with the given group appends to the receiver's
// existing groups. The keys of all subsequent attributes, whether added by With
// or in a Record, should be qualified by the sequence of group names.
func (l *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: l.handler.WithGroup(name), events: l.events}
}

// Handle looks to see if an event function needs to be executed for a given
// lol level and then formats its argument Record.
func (l *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if l.events.Debug != nil {
			l.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelInfo:
		if l.events.Info != nil {
			l.events.Info(ctx, toRecord(r))
		}

	case slog.LevelWarn:
		if l.events.Warn != nil {
			l.events.Warn(ctx, toRecord(r))
		}

	case slog.LevelError:
		if l.events.Error != nil {
			l.events.Error(ctx, toRecord(r))
		}
	}

	return l.handler.Handle(ctx, r)
}
