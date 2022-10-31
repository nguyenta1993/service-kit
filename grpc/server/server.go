package server

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	interceptors "github.com/gogovan-korea/ggx-kr-service-utils/grpc/interceptors"
	"github.com/gogovan-korea/ggx-kr-service-utils/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger             logger.Logger
	cfg                GrpcServerConfig
	GrpcServerInstance *grpc.Server
}

type GrpcServer interface {
	Run()
	Stop()
}

type GrpcServerConfig struct {
	Port              string
	Development       bool
	MaxConnectionIdle int
	Timeout           int
	MaxConnectionAge  int
	Time              int
}

func NewServer(logger logger.Logger, cfg GrpcServerConfig) (GrpcServer, *grpc.Server) {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Time) * time.Minute,
		}),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_recovery.UnaryServerInterceptor(),
				interceptors.Localizer(),
				interceptors.Logger(logger),
				otelgrpc.UnaryServerInterceptor(),
			),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor()),
	)

	grpc_prometheus.Register(grpcServer)

	if cfg.Development {
		reflection.Register(grpcServer)
	}

	return &Server{logger: logger, cfg: cfg, GrpcServerInstance: grpcServer}, grpcServer
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	lis, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
		panic(err)
	}

	go func() {
		s.logger.Info("GRPC server is listening at: ", zap.String("PORT", lis.Addr().String()))
		if err := s.GrpcServerInstance.Serve(lis); err != nil {
			s.logger.Fatal("failed to listen", zap.Error(err))
		}
	}()

	<-ctx.Done()

	_, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		fmt.Println("Close another connection")
		cancel()
	}()

	s.GrpcServerInstance.GracefulStop()
}

func (s *Server) Stop() {
	s.logger.Info("Stop GRPC server")
	s.GrpcServerInstance.GracefulStop()
}
