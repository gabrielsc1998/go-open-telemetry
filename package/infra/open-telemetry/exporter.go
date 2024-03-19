package opentelemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExporter(exporterType string, endpoint string, ctx context.Context) (sdktrace.SpanExporter, error) {
	switch exporterType {
	case "jaeger":
		return jaegerExporter(endpoint, ctx)
	case "zipkin":
		return zipkinExporter(endpoint, ctx)
	default:
		return nil, fmt.Errorf("exporter type %s not supported", exporterType)
	}
}

func jaegerExporter(endpoint string, ctx context.Context) (sdktrace.SpanExporter, error) {
	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}

func zipkinExporter(endpoint string, ctx context.Context) (sdktrace.SpanExporter, error) {
	traceExporter, err := zipkin.New(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}
