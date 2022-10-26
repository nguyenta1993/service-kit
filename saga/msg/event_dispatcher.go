package msg

import (
	"context"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
	"go.uber.org/zap"
)

// EventHandlerFunc function handlers for msg.Event
type EventHandlerFunc func(context.Context, Event) error

// EventDispatcher is a MessageReceiver for Events
type EventDispatcher struct {
	handlers map[string]EventHandlerFunc
	logger   logger.Logger
}

var _ MessageReceiver = (*EventDispatcher)(nil)

// NewEventDispatcher constructs a new EventDispatcher
func NewEventDispatcher(logger logger.Logger, options ...EventDispatcherOption) *EventDispatcher {
	c := &EventDispatcher{
		handlers: map[string]EventHandlerFunc{},
		logger:   logger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Info("msg.EventDispatcher constructed")

	return c
}

// Handle adds a new Event that will be handled by EventMessageFunc handler
func (d *EventDispatcher) Handle(evt core.Event, handler EventHandlerFunc) *EventDispatcher {
	d.logger.Info("event handler added", zap.String("EventName", evt.EventName()))
	d.handlers[evt.EventName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *EventDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	eventName, err := message.Headers().GetRequired(MessageEventName)
	if err != nil {
		d.logger.Error("error reading event name", zap.Error(err))
		return nil
	}

	logger := d.logger.With(
		zap.String("EventName", eventName),
		zap.String("MessageID", message.ID()),
	)

	logger.Debug("received event message")

	// check first for a handler of the event; It is possible events might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[eventName]
	if !exists {
		return nil
	}

	logger.Info("event handler found")

	event, err := core.DeserializeEvent(eventName, message.Payload())
	if err != nil {
		logger.Error("error decoding event message payload", zap.Error(err))
		return nil
	}

	evtMsg := eventMessage{event, message.Headers()}

	err = handler(ctx, evtMsg)
	if err != nil {
		logger.Error("event handler returned an error", zap.Error(err))
	}

	return err
}
