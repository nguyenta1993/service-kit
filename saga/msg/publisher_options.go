package msg

import "github.com/gogovan-korea/ggx-kr-service-utils/logger"

// PublisherOption options for PublisherPublisher
type PublisherOption func(*Publisher)

// WithPublisherLogger is an option to set the logger.Logger of the Publisher
func WithPublisherLogger(logger logger.Logger) PublisherOption {
	return func(publisher *Publisher) {
		publisher.logger = logger
	}
}
