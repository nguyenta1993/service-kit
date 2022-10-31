package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	j "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

type Config struct {
	ServiceName string
	HostPort    string
	Enable      bool
	LogSpans    bool
}

var tracer trace.Tracer

func tracerProvider(jaegerConfig Config) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	hostPort := strings.Split(jaegerConfig.HostPort, ":")
	exp, err := j.New(
		j.WithAgentEndpoint(
			j.WithAgentHost(hostPort[0]),
			j.WithAgentPort(hostPort[1]),
		),
	)

	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(jaegerConfig.ServiceName),
			attribute.String("serviceName", jaegerConfig.ServiceName),
		)),
	)
	tracer = tp.Tracer(jaegerConfig.ServiceName)
	return tp, nil
}
