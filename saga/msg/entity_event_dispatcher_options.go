package msg

import "github.com/gogovan-korea/ggx-kr-service-utils/logger"

// EntityEventDispatcherOption options for EntityEventDispatcher
type EntityEventDispatcherOption func(consumer *EntityEventDispatcher)

// WithEntityEventDispatcherLogger is an option to set the logger.Logger of the EntityEventDispatcher
func WithEntityEventDispatcherLogger(logger logger.Logger) EntityEventDispatcherOption {
	return func(dispatcher *EntityEventDispatcher) {
		dispatcher.logger = logger
	}
}
