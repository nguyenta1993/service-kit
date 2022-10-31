package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"strings"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
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
	if span := trace.SpanFromContext(ctx); span != nil {
		zLogger := otelzap.New(logger.WithOptions(zap.AddCallerSkip(1)).GetZapLogger())
		jaegerCtx := span.SpanContext()
		otelzap.ReplaceGlobals(zLogger)
		spanLogger := spanLogger{}
		spanLogger.spanFields = []zapcore.Field{
			zap.String("trace_id", jaegerCtx.TraceID().String()),
			zap.String("span_id", jaegerCtx.SpanID().String()),
		}
		return spanLogger
	}
	return logger
}
