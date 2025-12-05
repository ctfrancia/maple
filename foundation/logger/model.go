package logger

import (
	"context"
	"log/slog"
	"time"
)

type Level slog.Level

const (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
	LevelError Level = Level(slog.LevelError)
)

// Record represends the data that is being logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())

	f := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value.Any()
		return true
	}

	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: atts,
	}
}

// EventFn is a function that is called when an event is logged.
type EventFn func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
