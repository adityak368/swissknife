package tracing

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func Init(config OpenTracingConfig) opentracing.Tracer {
	// init open tracing
	udpTransport, _ := jaeger.NewUDPTransport(config.JaegerURL, 0)
	reporter := jaeger.NewRemoteReporter(udpTransport)
	sampler := jaeger.NewConstSampler(true)
	tracer, _ := jaeger.NewTracer(config.ServiceName, sampler, reporter)
	opentracing.SetGlobalTracer(tracer)
	return tracer
}

// StartSpanFromContext returns a new span with the given operation name and options. If a span
// is found in the context, it will be used as the parent of the resulting span.
func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// copy the metadata to prevent race
	md = metadata.Copy(md)

	// Find parent span.
	// First try to get span within current service boundary.
	// If there doesn't exist, try to get it from go-micro metadata(which is cross boundary)
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md)); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	sp := tracer.StartSpan(name, opts...)

	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return nil, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	return ctx, sp, nil
}

// EchoTracingMiddleware returns a middleware for echo
func EchoTracingMiddleware(ot opentracing.Tracer, spanName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if ot == nil {
				ot = opentracing.GlobalTracer()
			}
			name := fmt.Sprintf("%s: %s", spanName, c.Request().URL.Path)
			ctx, span, err := StartSpanFromContext(c.Request().Context(), ot, name)
			if err != nil {
				return err
			}
			defer span.Finish()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
