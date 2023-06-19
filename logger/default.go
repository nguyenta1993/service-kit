package logger

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type DefaultLogger struct {
	*zap.Logger
}

var defaultLogger Logger

func init() {
	defaultLogger = NewDefaultLogger(viper.GetString("logger.level"))
}

func NewDefaultLogger(consoleLevel string) Logger {
	if consoleLevel == "" ||
		(consoleLevel != DebugLevel && consoleLevel != InfoLevel && consoleLevel != WarnLevel && consoleLevel != ErrorLevel) {
		consoleLevel = InfoLevel
	}
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	logger := zap.New(
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), getZapLevel(consoleLevel)),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller())
	logger = logger.WithOptions(zap.AddCallerSkip(2))
	return zapLogger{logger}
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

func Debug(msg string, fields ...zapcore.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	defaultLogger.Panic(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	defaultLogger.Fatal(msg, fields...)
}

func Log(keyvals ...interface{}) error {
	return defaultLogger.Log(keyvals...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Info(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Debug(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Error(msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Fatal(msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithTrace(ctx, defaultLogger).Panic(msg, fields...)
}
