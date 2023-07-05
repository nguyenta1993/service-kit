package msg

import "github.com/nguyenta1993/service-kit/logger"

// CommandDispatcherOption options for CommandDispatcher
type CommandDispatcherOption func(consumer *CommandDispatcher)

// WithCommandDispatcherLogger is an option to set the logger.Logger of the CommandDispatcher
func WithCommandDispatcherLogger(logger logger.Logger) CommandDispatcherOption {
	return func(dispatcher *CommandDispatcher) {
		dispatcher.logger = logger
	}
}
