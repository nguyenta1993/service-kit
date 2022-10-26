package msg

import (
	"context"
	"strings"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"go.uber.org/zap"

	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
)

// CommandHandlerFunc function handlers for msg.Command
type CommandHandlerFunc func(context.Context, Command) ([]Reply, error)

// CommandDispatcher is a MessageReceiver for Commands
type CommandDispatcher struct {
	publisher ReplyMessagePublisher
	handlers  map[string]CommandHandlerFunc
	logger    logger.Logger
}

var _ MessageReceiver = (*CommandDispatcher)(nil)

// NewCommandDispatcher constructs a new CommandDispatcher
func NewCommandDispatcher(publisher ReplyMessagePublisher, logger logger.Logger, options ...CommandDispatcherOption) *CommandDispatcher {
	c := &CommandDispatcher{
		publisher: publisher,
		handlers:  map[string]CommandHandlerFunc{},
		logger:    logger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Info("msg.CommandDispatcher constructed")

	return c
}

// Handle adds a new Command that will be handled by handler
func (d *CommandDispatcher) Handle(cmd core.Command, handler CommandHandlerFunc) *CommandDispatcher {
	d.logger.Info("command handler added", zap.String("CommandName", cmd.CommandName()))
	d.handlers[cmd.CommandName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *CommandDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	commandName, err := message.Headers().GetRequired(MessageCommandName)
	if err != nil {
		d.logger.Error("error reading command name", zap.Error(err))
		return nil
	}

	logger := d.logger.With(
		zap.String("CommandName", commandName),
		zap.String("MessageID", message.ID()),
	)

	logger.Debug("received command message")

	// check first for a handler of the command; It is possible commands might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[commandName]
	if !exists {
		return nil
	}

	logger.Info("command handler found")

	command, err := core.DeserializeCommand(commandName, message.Payload())
	if err != nil {
		logger.Error("error decoding command message payload", zap.Error(err))
		return nil
	}

	replyChannel, err := message.Headers().GetRequired(MessageCommandReplyChannel)
	if err != nil {
		logger.Error("error reading reply channel", zap.Error(err))
		return nil
	}

	correlationHeaders := d.correlationHeaders(message.Headers())

	cmdMsg := commandMessage{command, correlationHeaders}

	replies, err := handler(ctx, cmdMsg)
	if err != nil {
		logger.Error("command handler returned an error", zap.Error(err))
		rerr := d.sendReplies(ctx, replyChannel, []Reply{WithFailure()}, correlationHeaders)
		if rerr != nil {
			logger.Error("error sending replies", zap.Error(rerr))
			return nil
		}
		return nil
	}

	err = d.sendReplies(ctx, replyChannel, replies, correlationHeaders)
	if err != nil {
		logger.Error("error sending replies", zap.Error(err))
		return nil
	}

	return nil
}

func (d *CommandDispatcher) sendReplies(ctx context.Context, replyChannel string, replies []Reply, correlationHeaders Headers) error {
	for _, reply := range replies {
		err := d.publisher.PublishReply(ctx, reply.Reply(),
			WithHeaders(reply.Headers()),
			WithHeaders(correlationHeaders),
			WithDestinationChannel(replyChannel),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *CommandDispatcher) correlationHeaders(headers Headers) Headers {
	replyHeaders := make(map[string]string)
	for key, value := range headers {
		if key == MessageCommandName {
			continue
		}

		if strings.HasPrefix(key, MessageCommandPrefix) {
			replyHeader := MessageReplyPrefix + key[len(MessageCommandPrefix):]
			replyHeaders[replyHeader] = value
		}
	}

	return replyHeaders
}
