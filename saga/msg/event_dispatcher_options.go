package msg

import "github.com/nguyenta1993/service-kit/logger"

// EventDispatcherOption options for EventDispatcher
type EventDispatcherOption func(consumer *EventDispatcher)

// WithEventDispatcherLogger is an option to set the logger.Logger of the EventDispatcher
func WithEventDispatcherLogger(logger logger.Logger) EventDispatcherOption {
	return func(dispatcher *EventDispatcher) {
		dispatcher.logger = logger
	}
}
