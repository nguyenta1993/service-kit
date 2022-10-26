package msg

import (
	"context"
	"sync"
	"time"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
	"go.uber.org/zap"
)

// CommandMessagePublisher interface
type CommandMessagePublisher interface {
	PublishCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error
}

// EntityEventMessagePublisher interface
type EntityEventMessagePublisher interface {
	PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error
}

// EventMessagePublisher interface
type EventMessagePublisher interface {
	PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error
}

// ReplyMessagePublisher interface
type ReplyMessagePublisher interface {
	PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error
}

// MessagePublisher interface
type MessagePublisher interface {
	Publish(ctx context.Context, message Message) error
}

var _ CommandMessagePublisher = (*Publisher)(nil)
var _ EntityEventMessagePublisher = (*Publisher)(nil)
var _ EventMessagePublisher = (*Publisher)(nil)
var _ ReplyMessagePublisher = (*Publisher)(nil)
var _ MessagePublisher = (*Publisher)(nil)

// Publisher send domain events, commands, and replies to the publisher
type Publisher struct {
	producer Producer
	logger   logger.Logger
	close    sync.Once
}

// NewPublisher constructs a new Publisher
func NewPublisher(producer Producer, logger logger.Logger, options ...PublisherOption) *Publisher {
	p := &Publisher{
		producer: producer,
		logger:   logger,
	}

	for _, option := range options {
		option(p)
	}

	p.logger.Info("msg.Publisher constructed")

	return p
}

// PublishCommand serializes a command into a message with command specific headers and publishes it to a producer
func (p *Publisher) PublishCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error {
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageCommandName:         command.CommandName(),
			MessageCommandReplyChannel: replyChannel,
		}),
	}

	if v, ok := command.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.With(
		zap.String("CommandName", command.CommandName()),
	)

	logger.Info("publishing command")

	payload, err := core.SerializeCommand(command)
	if err != nil {
		logger.Error("error serializing command payload", zap.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing command", zap.Error(err))
	}

	return err
}

// PublishReply serializes a reply into a message with reply specific headers and publishes it to a producer
func (p *Publisher) PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error {
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageReplyName: reply.ReplyName(),
		}),
	}

	if v, ok := reply.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.With(
		zap.String("ReplyName", reply.ReplyName()),
	)

	logger.Info("publishing reply")

	payload, err := core.SerializeReply(reply)
	if err != nil {
		logger.Error("error serializing reply payload", zap.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing reply", zap.Error(err))
	}

	return err
}

// PublishEntityEvents serializes entity events into messages with entity specific headers and publishes it to a producer
func (p *Publisher) PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error {
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageEventEntityID:   entity.ID(),
			MessageEventEntityName: entity.EntityName(),
			MessageChannel:         entity.EntityName(), // allow entity name and channel to overlap
		}),
	}

	if v, ok := entity.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	for _, event := range entity.Events() {
		logger := p.logger.With(
			zap.String("EntityID", entity.ID()),
			zap.String("EntityName", entity.EntityName()),
		)

		err := p.PublishEvent(ctx, event, msgOptions...)
		if err != nil {
			logger.Error("error publishing entity event", zap.Error(err))
			return err
		}
	}

	return nil
}

// PublishEvent serializes an event into a message with event specific headers and publishes it to a producer
func (p *Publisher) PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error {
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageEventName: event.EventName(),
		}),
	}

	if v, ok := event.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.With(
		zap.String("EventName", event.EventName()),
	)

	logger.Info("publishing event")

	payload, err := core.SerializeEvent(event)
	if err != nil {
		logger.Error("error serializing event payload", zap.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing event", zap.Error(err))
	}

	return err
}

// Publish sends a message off to a producer
func (p *Publisher) Publish(ctx context.Context, message Message) error {
	var err error
	var channel string

	channel, err = message.Headers().GetRequired(MessageChannel)
	if err != nil {
		return err
	}

	message.Headers()[Messagsagae] = time.Now().Format(time.RFC3339)

	// Published messages are request boundaries
	if id, exists := message.Headers()[MessageCorrelationID]; !exists || id == "" {
		message.Headers()[MessageCorrelationID] = core.GetCorrelationID(ctx)
	}

	if id, exists := message.Headers()[MessageCausationID]; !exists || id == "" {
		message.Headers()[MessageCausationID] = core.GetRequestID(ctx)
	}

	logger := p.logger.With(
		zap.String("MessageID", message.ID()),
		zap.String("CorrelationID", message.Headers()[MessageCorrelationID]),
		zap.String("CausationID", message.Headers()[MessageCausationID]),
		zap.String("Destination", channel),
		zap.Int("PayloadSize", len(message.Payload())),
	)

	logger.Info("publishing message")

	err = p.producer.Send(ctx, channel, message)
	if err != nil {
		logger.Error("error publishing message", zap.Error(err))
		return err
	}

	return nil
}

// Stop stops the publisher and underlying producer
func (p *Publisher) Stop(ctx context.Context) (err error) {
	defer p.logger.Info("publisher stopped")
	p.close.Do(func() {
		err = p.producer.Close(ctx)
	})

	return
}
