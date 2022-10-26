package interceptors

import (
	"context"

	"github.com/tikivn/s14e-backend-utils/localization"

	"google.golang.org/grpc"
)

func Localizer() func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		localization.NewLocalizer(localization.ResourceConfig{
			Lang:   "vi",
			Accept: "vi",
		})

		reply, err := handler(ctx, req)

		return reply, err
	}
}
