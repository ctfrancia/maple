package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const defaultTraceID = "0000000000000000"

// Config defines the information needed to create a new tracer.
type Config struct {
	ServiceName    string
	Host           string
	ExcludedRoutes map[string]struct{}
	Probability    float64
}

func InitTracing(log *logger.Logger, cfg Config) (trace.TracerProvider, func(ctx context.Context), error) {
	// will need to add my own settings that are not default
	// will need to go look at opentelemetry docs

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(), // this should be configurable
			otlptracegrpc.WithEndpoint(cfg.Host),
		),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("creating new exporter: %w", err)
	}

	var tracerProvider trace.TracerProvider
	teardown := func(ctx context.Context) {}

	switch cfg.Host {
	case "":
		log.Info(context.Background(), "OTEL", "tracer", "NOOP")
		tracerProvider = noop.NewTracerProvider()

	default:
		log.Info(context.Background(), "OTEL", "tracer", cfg.Host)

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.ParentBased(newEndpointExcluder(cfg.ExcludedRoutes, cfg.Probability))),
			sdktrace.WithBatcher(exporter,
				sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
				sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			),
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(cfg.ServiceName),
				),
			),
		)

		teardown = func(ctx context.Context) {
			tp.Shutdown(ctx)
		}

		tracerProvider = tp

	}

	// We must set this provider as the flobal provider for things to work,
	// but we pass this provider areound the program where needed to collect our
	// traces.
	otel.SetTracerProvider(tracerProvider)

	// Extract incoming trace contexts and the headers we set in ongoing requests
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider, teardown, nil
}

// InjectTracing initializes the request for tracing by writing otel related
// information into the response and saving the tracer and trace id in the context
// for later use.
func InjectTracing(ctx context.Context, tracer trace.Tracer) context.Context {
	ctx = setTracer(ctx, tracer)

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	if traceID == defaultTraceID {
		traceID = uuid.NewString()
	}

	ctx = setTraceID(ctx, traceID)

	return ctx
}

// AddSpan adds an otel span to the existing trace.
func AddSpan(ctx context.Context, spanName string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer, ok := ctx.Value(tracerKey).(trace.Tracer)
	if !ok || tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	ctx, span := tracer.Start(ctx, spanName)

	span.SetAttributes(keyValues...)

	return ctx, span
}
