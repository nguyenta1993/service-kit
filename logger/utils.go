package logger

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func WithSpan(logger Logger, ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		spanLogger := spanLogger{span: span, logger: logger.WithOptions(zap.AddCallerSkip(1))}

		if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
			spanLogger.spanFields = []zapcore.Field{
				zap.String("trace_id", jaegerCtx.TraceID().String()),
				zap.String("span_id", jaegerCtx.SpanID().String()),
			}
		}

		return spanLogger
	}
	return logger
}
