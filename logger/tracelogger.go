package logger

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type traceLogger struct {
	otelzap.LoggerWithCtx
}

func (t traceLogger) Log(keyvals ...interface{}) error {
	t.LoggerWithCtx.WithOptions(zap.AddCallerSkip(2)).Sugar().Infow("log", keyvals...)
	return nil
}

func (t traceLogger) With(fields ...zapcore.Field) Logger {
	t.LoggerWithCtx.ZapLogger().With(fields...)
	return t
}

func (t traceLogger) WithOptions(opts ...zap.Option) Logger {
	t.LoggerWithCtx = t.LoggerWithCtx.WithOptions(opts...)
	return t
}

func (t traceLogger) GetZapLogger() *zap.Logger {
	return t.ZapLogger()
}

// WithTrace  use logger with tracing context
func WithTrace(logger Logger, ctx context.Context) Logger {
	log := otelzap.New(logger.GetZapLogger()).Ctx(ctx)
	return traceLogger{log}
}
