// logger provides a zap logger that is used for logging
package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ctfrancia/maple/internal/core/ports"
)

// ZapLogger is a logger that uses the Zap library
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger
func NewZapLogger(env string) ports.Logger {
	var config zap.Config

	if env == "prod" {
		config = zap.NewProductionConfig()
		// Add log rotation and other production configs here
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return &ZapLogger{
		logger: logger,
	}
}

// convertFields converts domain LogFields to zap.Fields
func (z *ZapLogger) convertFields(fields []ports.LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		switch v := field.Value.(type) {
		case error:
			zapFields[i] = zap.Error(v)
		case string:
			zapFields[i] = zap.String(field.Key, v)
		case int:
			zapFields[i] = zap.Int(field.Key, v)
		case int64:
			zapFields[i] = zap.Int64(field.Key, v)
		case float64:
			zapFields[i] = zap.Float64(field.Key, v)
		case bool:
			zapFields[i] = zap.Bool(field.Key, v)
		default:
			zapFields[i] = zap.Any(field.Key, v)
		}
	}
	return zapFields
}

// Debug logs a message at DebugLevel
func (z *ZapLogger) Debug(ctx context.Context, msg string, fields ...ports.LogField) {
	z.logger.Debug(msg, z.convertFields(fields)...)
}

// Info logs a message at InfoLevel
func (z *ZapLogger) Info(ctx context.Context, msg string, fields ...ports.LogField) {
	z.logger.Info(msg, z.convertFields(fields)...)
}

// warn logs a message at WarnLevel
func (z *ZapLogger) Warn(ctx context.Context, msg string, fields ...ports.LogField) {
	z.logger.Warn(msg, z.convertFields(fields)...)
}

// Error logs a message at ErrorLevel
func (z *ZapLogger) Error(ctx context.Context, msg string, fields ...ports.LogField) {
	z.logger.Error(msg, z.convertFields(fields)...)
}

// Fatal logs a message at FatalLevel
func (z *ZapLogger) Fatal(ctx context.Context, msg string, fields ...ports.LogField) {
	z.logger.Fatal(msg, z.convertFields(fields)...)
}
