package middleware

import (
    "go.opentelemetry.io/otel/propagation"
    oteltrace "go.opentelemetry.io/otel/trace"
)

type config struct {
    TracerProvider oteltrace.TracerProvider
    Propagators    propagation.TextMapPropagator
}

// Option specifies instrumentation configuration options.
type Option interface {
    apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
    o(c)
}
