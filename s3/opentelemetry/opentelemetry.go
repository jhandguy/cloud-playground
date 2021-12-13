package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var provider *sdktrace.TracerProvider

func StartTracing(ctx context.Context, endpoint string) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("s3"),
		),
	)
	if err != nil {
		return err
	}

	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return err
	}

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exp),
	)
	otel.SetTracerProvider(provider)
	return nil
}

func StopTracing(ctx context.Context) error {
	if provider == nil {
		return nil
	}
	return provider.Shutdown(ctx)
}

func GetTracer(name string) trace.Tracer {
	if provider == nil {
		return nil
	}
	return provider.Tracer(name)
}

func GetTraceID(ctx context.Context) string {
	if provider == nil {
		return ""
	}
	sp := trace.SpanFromContext(ctx)
	return sp.SpanContext().TraceID().String()
}
