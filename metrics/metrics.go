package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/gogovan/ggx-kr-service-utils/http/middlewares"
	"github.com/gogovan/ggx-kr-service-utils/logger"
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
		metricsServer.Use(middlewares.Logging(logger.GetDefaultLogger()))
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
