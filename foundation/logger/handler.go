package logger

import "log/slog"

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
