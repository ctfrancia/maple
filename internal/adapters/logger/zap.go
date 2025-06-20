package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type ZapAdapter struct {
	logger *zap.Logger
}

// NewZapLogger creates a new zap logger based on environment
func NewZapLogger(env string) *ZapAdapter {
	zapLogger := newZapLogger(env)
	return &ZapAdapter{logger: zapLogger}
}

// Your existing code slightly modified
func newZapLogger(env string) *zap.Logger {
	var core zapcore.Core
	err := os.MkdirAll("./internal/logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	isDev := env == "dev" || env == "test"
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logLevel := zap.InfoLevel
	if isDev {
		// In development: Use a console encoder and write to stderr
		logLevel = zap.DebugLevel
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stderr),
			logLevel,
		)
	} else {
		// In production/other environments: Use JSON encoder and log rotation
		logRotator := &lumberjack.Logger{
			Filename:   "./internal/logs/buho.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(logRotator),
			logLevel,
		)
	}

	return zap.New(core, zap.AddCaller())
}

// Implement the Logger interface
func (z *ZapAdapter) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *ZapAdapter) Info(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapAdapter) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *ZapAdapter) Error(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *ZapAdapter) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

/*
// Optional: Add WithContext for tracing/request IDs
func (z *ZapAdapter) WithContext(ctx context.Context) *zap.Logger {
	// You could extract values from context here
	// For example: request ID, user ID, etc.
	return z
}
*/
