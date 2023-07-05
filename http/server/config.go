package server

import (
	"github.com/nguyenta1993/service-kit/logger"
	"github.com/nguyenta1993/service-kit/metrics"
)

type HttpServerConfig struct {
	Port            string
	Development     bool
	ShutdownTimeout int

	Resources    []string
	RateLimiting *RateLimitingConfig
	Name         string
	AllowOrigins []string
	MetricConfig *metrics.MetricsConfig
}

type RateLimitingConfig struct {
	RateFormat string
}

type HttpServerOption func(*Server)

func WithLogger(log logger.Logger) HttpServerOption {
	return func(s *Server) {
		s.log = log
	}
}
