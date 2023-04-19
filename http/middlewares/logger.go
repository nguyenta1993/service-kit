package middlewares

import (
	"bytes"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
	"io"
	"time"

	"github.com/gogovan/ggx-kr-service-utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var body []byte
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ = io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		c.Next()
		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if len(c.Request.URL.RawQuery) > 0 {
			fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
		}

		if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		if span := trace.SpanFromContext(c.Request.Context()).SpanContext(); span.IsValid() {
			fields = append(fields, zap.String("trace_id", span.TraceID().String()))
		}

		if len(body) > 0 {
			fields = append(fields, zap.String("body", string(body)))
		}
		fields = append(fields, zap.Duration("latency", time.Since(start)))
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info(c.Request.URL.Path, fields...)
		}
	}
}

func Logging(logger logger.Logger, skips ...string) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger.GetZapLogger(), &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		SkipPaths:  skips,
		Context: func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}
			// log trace and span ID
			if span := trace.SpanFromContext(c.Request.Context()).SpanContext(); span.IsValid() {
				fields = append(fields, zap.String("trace_id", span.TraceID().String()))
				fields = append(fields, zap.String("span_id", span.SpanID().String()))
			}
			return fields
		},
	})
}

func Recovery(logger logger.Logger) gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger.GetZapLogger(), true)
}

func Tracing(name string) gin.HandlerFunc {
	return otelgin.Middleware(name, otelgin.WithPropagators(otel.GetTextMapPropagator()))
}

func RequestId() gin.HandlerFunc {
	return requestid.New()
}

func Gzip() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}
