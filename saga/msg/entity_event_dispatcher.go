package msg

import (
	"context"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
	"go.uber.org/zap"
)

// EntityEventHandlerFunc function handlers for msg.EntityEvent
type EntityEventHandlerFunc func(context.Context, EntityEvent) error

// EntityEventDispatcher is a MessageReceiver for DomainEvents
type EntityEventDispatcher struct {
	handlers map[string]EntityEventHandlerFunc
	logger   logger.Logger
}

var _ MessageReceiver = (*EntityEventDispatcher)(nil)

// NewEntityEventDispatcher constructs a new EntityEventDispatcher
func NewEntityEventDispatcher(logger logger.Logger, options ...EntityEventDispatcherOption) *EntityEventDispatcher {
	c := &EntityEventDispatcher{
		handlers: map[string]EntityEventHandlerFunc{},
		logger:   logger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Info("msg.EntityEventDispatcher constructed")

	return c
}

// Handle adds a new Event that will be handled by EventMessageFunc handler
func (d *EntityEventDispatcher) Handle(evt core.Event, handler EntityEventHandlerFunc) *EntityEventDispatcher {
	d.logger.Info("entity event handler added", zap.String("EventName", evt.EventName()))
	d.handlers[evt.EventName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *EntityEventDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	eventName, err := message.Headers().GetRequired(MessageEventName)
	if err != nil {
		d.logger.Error("error reading event name", zap.Error(err))
		return nil
	}

	entityName, err := message.Headers().GetRequired(MessageEventEntityName)
	if err != nil {
		d.logger.Error("error reading entity name", zap.Error(err))
		return nil
	}

	entityID, err := message.Headers().GetRequired(MessageEventEntityID)
	if err != nil {
		d.logger.Error("error reading entity id", zap.Error(err))
		return nil
	}

	logger := d.logger.With(
		zap.String("EntityName", entityName),
		zap.String("EntityID", entityID),
		zap.String("EventName", eventName),
		zap.String("MessageID", message.ID()),
	)

	logger.Debug("received entity event message")

	// check first for a handler of the event; It is possible events might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[eventName]
	if !exists {
		return nil
	}

	logger.Info("entity event handler found")

	event, err := core.DeserializeEvent(eventName, message.Payload())
	if err != nil {
		logger.Error("error decoding entity event message payload", zap.Error(err))
		return nil
	}

	evtMsg := entityEventMessage{entityID, entityName, event, message.Headers()}

	err = handler(ctx, evtMsg)
	if err != nil {
		logger.Error("entity event handler returned an error", zap.Error(err))
	}

	return err
}
