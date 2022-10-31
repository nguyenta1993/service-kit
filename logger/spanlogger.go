package logger

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger     Logger
	span       trace.Span
	spanFields []zapcore.Field
}

func (sl spanLogger) Debug(msg string, fields ...zapcore.Field) {
	sl.logger.Debug(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Info(msg string, fields ...zapcore.Field) {
	sl.logToSpan("info", msg, fields...)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Warn(msg string, fields ...zapcore.Field) {
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Error(msg string, fields ...zapcore.Field) {
	sl.logToSpan("error", msg, fields...)
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	sl.logToSpan("fatal", msg, fields...)
	sl.logger.Fatal(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Panic(msg string, fields ...zapcore.Field) {
	sl.logToSpan("panic", msg, fields...)
	sl.logger.Fatal(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Log(keyvals ...interface{}) error {
	return sl.logger.Log(keyvals...)
}

func (sl spanLogger) With(fields ...zapcore.Field) Logger {
	return spanLogger{logger: sl.logger.With(fields...), span: sl.span, spanFields: sl.spanFields}
}

func (sl spanLogger) WithOptions(opts ...zap.Option) Logger {
	return spanLogger{logger: sl.logger.WithOptions(opts...), span: sl.span, spanFields: sl.spanFields}
}

func (l spanLogger) GetZapLogger() *zap.Logger {
	return l.logger.GetZapLogger()
}

func (sl spanLogger) logToSpan(level string, msg string, fields ...zapcore.Field) {
	fa := fieldAdapter(make([]attribute.KeyValue, 0, 2+len(fields)))
	fa = append(fa, attribute.String("event", msg))
	fa = append(fa, attribute.String("level", level))
	for _, field := range fields {
		field.AddTo(&fa)
	}
	sl.span.SetAttributes(fa...)
}

type fieldAdapter []attribute.KeyValue

func (fa *fieldAdapter) AddBool(key string, value bool) {
	*fa = append(*fa, attribute.Bool(key, value))
}

func (fa *fieldAdapter) AddFloat64(key string, value float64) {
	*fa = append(*fa, attribute.Float64(key, value))
}

func (fa *fieldAdapter) AddFloat32(key string, value float32) {
	*fa = append(*fa, attribute.Float64(key, float64(value)))
}

func (fa *fieldAdapter) AddInt(key string, value int) {
	*fa = append(*fa, attribute.Int(key, value))
}

func (fa *fieldAdapter) AddInt64(key string, value int64) {
	*fa = append(*fa, attribute.Int64(key, value))
}

func (fa *fieldAdapter) AddInt32(key string, value int32) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt16(key string, value int16) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt8(key string, value int8) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint(key string, value uint) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint64(key string, value uint64) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint32(key string, value uint32) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint16(key string, value uint16) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint8(key string, value uint8) {
	*fa = append(*fa, attribute.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUintptr(key string, value uintptr)                        {}
func (fa *fieldAdapter) AddArray(key string, marshaler zapcore.ArrayMarshaler) error { return nil }
func (fa *fieldAdapter) AddComplex128(key string, value complex128)                  {}
func (fa *fieldAdapter) AddComplex64(key string, value complex64)                    {}
func (fa *fieldAdapter) AddObject(key string, value zapcore.ObjectMarshaler) error   { return nil }
func (fa *fieldAdapter) AddReflected(key string, value interface{}) error            { return nil }
func (fa *fieldAdapter) OpenNamespace(key string)                                    {}

func (fa *fieldAdapter) AddDuration(key string, value time.Duration) {
	// TODO inefficient
	*fa = append(*fa, attribute.String(key, value.String()))
}

func (fa *fieldAdapter) AddTime(key string, value time.Time) {
	// TODO inefficient
	*fa = append(*fa, attribute.String(key, value.String()))
}

func (fa *fieldAdapter) AddBinary(key string, value []byte) {}

func (fa *fieldAdapter) AddByteString(key string, value []byte) {}

func (fa *fieldAdapter) AddString(key, value string) {
	if key != "" && value != "" {
		*fa = append(*fa, attribute.String(key, value))
	}
}
