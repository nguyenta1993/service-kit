package interceptors

import (
	"context"

	"github.com/gogovan-korea/ggx-kr-service-utils/localization"

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
