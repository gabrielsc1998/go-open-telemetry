package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
