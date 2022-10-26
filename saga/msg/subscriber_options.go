package msg

import "github.com/gogovan-korea/ggx-kr-service-utils/logger"

// SubscriberOption options for MessageConsumers
type SubscriberOption func(*Subscriber)

// WithSubscriberLogger is an option to set the logger.Logger of the Subscriber
func WithSubscriberLogger(logger logger.Logger) SubscriberOption {
	return func(subscriber *Subscriber) {
		subscriber.logger = logger
	}
}
