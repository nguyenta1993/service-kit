package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenta1993/service-kit/http/middlewares"
	"github.com/nguyenta1993/service-kit/logger"
	"go.uber.org/zap"
)

type MetricsConfig struct {
	PrometheusPort string `mapstructure:"prometheusPort"`
	PrometheusPath string `mapstructure:"prometheusPath"`
}

var m *Monitor

func Run(cfg *MetricsConfig) func() {
	return func() {
		gin.SetMode(gin.ReleaseMode)
		metricsServer := gin.New()
		metricsServer.Use(middlewares.Logging())
		m.SetMetricPath(cfg.PrometheusPath)
		m.Expose(metricsServer)
		logger.Info("Metrics server is running on port", zap.String("Metrics port", cfg.PrometheusPort))
		if err := metricsServer.Run(cfg.PrometheusPort); err != nil {
			logger.Error("metricsServer.Run", zap.Error(err))
		}
	}
}

func init() {
	m = GetMonitor()
}

func Use(router gin.IRoutes) {
	m.UseWithoutExposingEndpoint(router)
}
