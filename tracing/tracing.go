package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func UseOpenTelemetry(jaegerConfig Config) {
	tracerProvider, err := tracerProvider(jaegerConfig)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
