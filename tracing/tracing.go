package tracing

import (
	"github.com/opentracing/opentracing-go"
)

func UseTracing(jaegerConfig Config) {
	tracer, _, err := NewJaegerTracer(jaegerConfig)
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
}
