package trace

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"os"
	"time"
)

func Tracer() trace.Tracer {
	return otel.Tracer("")
}

func End(ctx context.Context, span trace.Span, err error) {
	// 判断是否有异常发生，如果有则设置一些异常信息
	if err != nil {
		// 记录异常
		span.RecordError(err)
		// 设置span 属性
		span.SetAttributes(
			// 设置事件为异常
			attribute.String("event", "error"),
			// 设置 message 为 err.Error().
			attribute.String("message", err.Error()),
		)
		//设置了 span 的状态
		span.SetStatus(codes.Error, err.Error())
	} else {
		// 如果没有发生异常，span 状态则为 ok
		span.SetStatus(codes.Ok, "OK")
	}
	// 中止 span
	span.End()
}

type CleanFunc func(ctx context.Context, tp *tracesdk.TracerProvider)

func NewTracerProvider(url, service string) (*tracesdk.TracerProvider, CleanFunc, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, nil, err
	}

	host, _ := os.Hostname()
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("hostname", host),
			semconv.ServiceNameKey.String(service), // 服务名
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Cleanly shutdown and flush telemetry when the application exits.
	cleanup := func(ctx context.Context, tp *tracesdk.TracerProvider) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("tracerProvider shutdown error")
		}
		log.Debug().Msg("tracerProvider shutdown")
	}

	return tp, cleanup, nil
}

// Ideally, when tracing application code, spans are created and managed in the application framework.
// You can access the current active span via the context object.
// which can be used for adding application specific attributes and events.
func printSpan(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	fmt.Printf("current span: %v\n", span)
}

// AddAttribute
// Example of adding an attribute to a span
// 向Span中添加属性值。
// To avoid collisions, always namespace your attribute keys using dot notation.
// like this: span.SetAttributes(attribute.Int("project.id", 2));
func addAttribute(ctx context.Context, key, value string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String(key, value))
}

// example of adding an event to a span
// 向Span中添加事件。
// Span events are a form of structured logging.
// Each event has a name, a timestamp, and a set of attributes.
// When events are added to a span, they inherit the span's context.
// This additional context allows events to be searched, filtered, and grouped by trace ID and other span attributes.
func addEvent(ctx context.Context, eventName string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("event1", trace.WithAttributes(
		attribute.String("event-attr1", "event-string1"),
		attribute.Int64("event-attr2", 10)))
}

// example of recording an exception
// 记录Span结果以及错误信息。
// Marking the span as an error is independent from recordings exceptions.
// To mark the entire span as an error, and have it count against error rates,
// set the SpanStatus to any value other than OK.
func recordException(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(errors.New("exception has occurred"))
	span.SetStatus(codes.Error, "internal error")
}
