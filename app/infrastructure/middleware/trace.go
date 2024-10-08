package middleware

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/propagation"
    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
    oteltrace "go.opentelemetry.io/otel/trace"
)

const (
    tracerKey  = "otel-go-tracer"
    tracerName = "otelgin"
)

// TraceMiddleware returns middleware that will trace incoming requests.
// The service parameter should describe the name of the (virtual)
// server handling the request.
func TraceMiddleware(service string, opts ...Option) gin.HandlerFunc {
    cfg := config{}
    for _, opt := range opts {
        opt.apply(&cfg)
    }
    if cfg.TracerProvider == nil {
        cfg.TracerProvider = otel.GetTracerProvider()
    }
    tracer := cfg.TracerProvider.Tracer(
        tracerName,
        oteltrace.WithInstrumentationVersion(otelgin.SemVersion()),
    )
    if cfg.Propagators == nil {
        cfg.Propagators = otel.GetTextMapPropagator()
    }
    return func(c *gin.Context) {
        c.Set(tracerKey, tracer)
        savedCtx := c.Request.Context()
        defer func() {
            c.Request = c.Request.WithContext(savedCtx)
        }()
        ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
        opts := []oteltrace.SpanStartOption{
            oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
            oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
            oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.FullPath(), c.Request)...),
            oteltrace.WithSpanKind(oteltrace.SpanKindServer),
        }
        spanName := c.FullPath()
        if spanName == "" {
            spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
        }
        ctx, span := tracer.Start(ctx, spanName, opts...)
        defer span.End()

        // pass the span through the request context
        c.Request = c.Request.WithContext(ctx)

        // serve the request to the next middleware
        c.Next()

        status := c.Writer.Status()
        attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
        spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
        span.SetAttributes(attrs...)
        span.SetStatus(spanStatus, spanMessage)
        if len(c.Errors) > 0 {
            span.SetStatus(codes.Error, c.Errors.String())
            span.RecordError(c.Errors[0])
        }
    }
}
