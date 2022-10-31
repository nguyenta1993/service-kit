package tracing

import (
	"go.opentelemetry.io/otel"
)

func UseOpenTelemetry(jaegerConfig Config) {
	tracerProvider, err := tracerProvider(jaegerConfig)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)
}
