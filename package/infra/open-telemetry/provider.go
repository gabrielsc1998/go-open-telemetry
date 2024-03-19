package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Provider struct {
	ServiceName       string
	CollectorEndpoint string
}

func NewProvider(serviceName, collectorEndpoint string) *Provider {
	return &Provider{
		ServiceName:       serviceName,
		CollectorEndpoint: collectorEndpoint,
	}
}

func (p *Provider) Start() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(p.ServiceName),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	_, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	traceExporter, err := NewExporter(
		"zipkin",
		p.CollectorEndpoint,
		ctx,
	)
	// traceExporter, err := NewExporter(
	// 	"jaeger",
	// 	p.CollectorEndpoint,
	// 	ctx,
	// )

	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceExporter.Shutdown, nil
}
