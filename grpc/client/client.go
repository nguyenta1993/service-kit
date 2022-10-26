package client

import (
	"context"
	"crypto/tls"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/tikivn/s14e-backend-utils/grpc/interceptors"
	"github.com/tikivn/s14e-backend-utils/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewClientConn(ctx context.Context, logger logger.Logger, port string, development bool) *grpc.ClientConn {
	unaryInterceptorOption := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
				grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
				grpc_retry.WithMax(backoffRetries),
			),
			interceptors.ClientLogger(logger),
		))

	var opts grpc.DialOption
	if development {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		opts = grpc.WithTransportCredentials(creds)
	} else {
		opts = grpc.WithInsecure()
	}

	clientConn, err := grpc.DialContext(
		ctx,
		port,
		unaryInterceptorOption,
		opts,
	)

	if err != nil {
		logger.Error("NewClientConn", zap.Error(err))
		return nil
	}

	return clientConn
}
