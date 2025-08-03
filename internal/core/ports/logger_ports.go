// Package ports contains the ports for the Maple
package ports

import "context"

// LogField represents a structured logging field
type LogField struct {
	Key   string
	Value any
}

// Logger defines the logging contract for the domain
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...LogField)
	Info(ctx context.Context, msg string, fields ...LogField)
	Warn(ctx context.Context, msg string, fields ...LogField)
	Error(ctx context.Context, msg string, fields ...LogField)
	Fatal(ctx context.Context, msg string, fields ...LogField)
}

// Helper functions to create log fields
func String(key, value string) LogField {
	return LogField{Key: key, Value: value}
}

func Int(key string, value int) LogField {
	return LogField{Key: key, Value: value}
}

func Error(key string, err error) LogField {
	return LogField{Key: key, Value: err}
}

func Any(key string, value any) LogField {
	return LogField{Key: key, Value: value}
}
