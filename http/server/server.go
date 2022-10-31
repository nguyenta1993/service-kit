package server

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	m "github.com/gogovan-korea/ggx-kr-service-utils/http/middlewares"
	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"go.uber.org/zap"
)

type Server struct {
	logger             logger.Logger
	cfg                HttpServerConfig
	Router             *gin.Engine
	httpServerInstance *http.Server
}

type HttpServer interface {
	Run()
	Stop()
}

func NewServer(logger logger.Logger, cfg HttpServerConfig) (HttpServer, *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	httpServerInstance := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	router.Use(m.SetLanguage(cfg.Resources))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(m.LoggerMiddleware(logger))
	router.Use(otelgin.Middleware(cfg.Name))
	if cfg.RateLimiting != nil {
		router.Use(m.RateLimittingMiddleware(logger, router, cfg.RateLimiting.RateFormat))
	}

	return &Server{logger: logger, cfg: cfg, Router: router, httpServerInstance: httpServerInstance}, router
}

func (s *Server) Run() {
	go func() {
		s.logger.Info("Http server is listening at: ", zap.String("PORT", s.cfg.Port))
		if err := s.httpServerInstance.ListenAndServe(); err != nil {
			s.logger.Error("failed to listen", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	s.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.httpServerInstance.Shutdown(ctx); err != nil {
		s.logger.Error("Server Shutdown:", zap.Error(err))
	}

	close(quit)
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := s.httpServerInstance.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server Shutdown:", zap.Error(err))
	}
}
