package healthcheck

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/gogovan/ggx-kr-service-utils/constants"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gogovan/ggx-kr-service-utils/logger"
	"github.com/heptiolabs/healthcheck"
	"go.uber.org/zap"
)

func RunHealthCheck(
	ctx context.Context,
	logger logger.Logger,
	cfg *HealthcheckConfig,
	readDb *sqlx.DB,
	writeDb *sqlx.DB,
	redis redis.UniversalClient,
	client *kafka.Conn,
) func() {
	return func() {
		health := healthcheck.NewHandler()

		livenessCheck(ctx, cfg.GoroutineThreshold, health)
		itv := time.Duration(cfg.Interval) * time.Second
		readinessCheck(ctx, logger, health, itv, readDb, writeDb, redis, client)

		logMd := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Our middleware logic goes here...
				rw := NewResponseWriter(w)
				next.ServeHTTP(rw, r)
				statusCode := rw.Code()
				logger.Info(
					"Response information",
					zap.String("status_code", strconv.Itoa(statusCode)),
					zap.String("Method", r.Method),
					zap.String("URL", r.RequestURI),
				)
			})
		}

		logger.Info("Heathcheck server listening on port", zap.String("Port", cfg.Port))
		if err := http.ListenAndServe(cfg.Port, logMd(health)); err != nil {
			logger.Warn("Heathcheck server", zap.Error(err))
		}
	}
}

func livenessCheck(ctx context.Context, goRoutinesThreshold int, health healthcheck.Handler) {
	health.AddLivenessCheck(constants.GoroutineThreshold, healthcheck.GoroutineCountCheck(goRoutinesThreshold))
}

func readinessCheck(
	ctx context.Context,
	logger logger.Logger,
	health healthcheck.Handler,
	interval time.Duration,
	readDb *sqlx.DB,
	writeDb *sqlx.DB,
	redis redis.UniversalClient,
	client *kafka.Conn,
) {

	if readDb != nil {
		health.AddReadinessCheck(constants.ReadDatabase, healthcheck.AsyncWithContext(ctx, func() (err error) {
			err = readDb.DB.PingContext(ctx)
			if err != nil {
				logger.Error("Read database", zap.Error(err))
			}
			return
		}, interval))
	}
	if writeDb != nil {
		health.AddReadinessCheck(constants.WriteDatabase, healthcheck.AsyncWithContext(ctx, func() (err error) {
			err = writeDb.DB.PingContext(ctx)
			if err != nil {
				logger.Error("Readiness check write database", zap.Error(err))
			}
			return
		}, interval))
	}

	if redis != nil {
		health.AddReadinessCheck(constants.Redis, healthcheck.AsyncWithContext(ctx, func() error {
			err := redis.Ping(ctx).Err()
			if err != nil {
				logger.Error("Redis Readiness Check Fail", zap.Error(err))
			}
			return err
		}, interval))
	}

	if client != nil {
		health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
			_, err := client.Brokers()
			if err != nil {
				logger.Error("Kafka Readiness Check Fail", zap.Error(err))
			}
			return err
		}, interval))
	}
}

type ResponseWriter struct {
	http.ResponseWriter

	code int
	size int
}

// Returns a new `ResponseWriter` type by decorating `http.ResponseWriter` type.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

// Overrides `http.ResponseWriter` type.
func (r *ResponseWriter) WriteHeader(code int) {
	if r.Code() == 0 {
		r.code = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// Overrides `http.ResponseWriter` type.
func (r *ResponseWriter) Write(body []byte) (int, error) {
	if r.Code() == 0 {
		r.WriteHeader(http.StatusOK)
	}

	var err error
	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

// Overrides `http.Flusher` type.
func (r *ResponseWriter) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		if r.Code() == 0 {
			r.WriteHeader(http.StatusOK)
		}

		fl.Flush()
	}
}

// Overrides `http.Hijacker` type.
func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()
}

// Returns response status code.
func (r *ResponseWriter) Code() int {
	return r.code
}

// Returns response size.
func (r *ResponseWriter) Size() int {
	return r.size
}
