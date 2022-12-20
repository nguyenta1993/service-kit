package saga

import (
	"context"
	"strings"

	"github.com/gogovan/ggx-kr-service-utils/logger"
	"github.com/gogovan/ggx-kr-service-utils/saga/core"
	"github.com/gogovan/ggx-kr-service-utils/saga/msg"
	"go.uber.org/zap"
)

// CommandHandlerFunc function handlers for saga.Command
type CommandHandlerFunc func(context.Context, Command) ([]msg.Reply, error)

// CommandDispatcher is a MessageReceiver for Commands
type CommandDispatcher struct {
	publisher msg.ReplyMessagePublisher
	handlers  map[string]CommandHandlerFunc
	logger    logger.Logger
}

var _ msg.MessageReceiver = (*CommandDispatcher)(nil)

// NewCommandDispatcher constructs a new CommandDispatcher
func NewCommandDispatcher(logger logger.Logger, publisher msg.ReplyMessagePublisher, options ...CommandDispatcherOption) *CommandDispatcher {
	c := &CommandDispatcher{
		publisher: publisher,
		handlers:  map[string]CommandHandlerFunc{},
		logger:    logger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Info("saga.CommandDispatcher constructed")

	return c
}

// Handle adds a new Command that will be handled by handler
func (d *CommandDispatcher) Handle(cmd core.Command, handler CommandHandlerFunc) *CommandDispatcher {
	d.logger.Info("saga command handler added", zap.String("CommandName", cmd.CommandName()))
	d.handlers[cmd.CommandName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *CommandDispatcher) ReceiveMessage(ctx context.Context, message msg.Message) error {
	commandName, sagaID, sagaName, err := d.commandMessageInfo(message)
	if err != nil {
		return nil
	}

	logger := d.logger.With(
		zap.String("CommandName", commandName),
		zap.String("SagaName", sagaName),
		zap.String("SagaID", sagaID),
		zap.String("MessageID", message.ID()),
	)

	logger.Debug("received saga command message")

	// check first for a handler of the command; It is possible commands might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[commandName]
	if !exists {
		return nil
	}

	logger.Info("saga command handler found")

	command, err := core.DeserializeCommand(commandName, message.Payload())
	if err != nil {
		logger.Error("error decoding saga command message payload", zap.Error(err))
		return nil
	}

	replyChannel, err := message.Headers().GetRequired(msg.MessageCommandReplyChannel)
	if err != nil {
		logger.Error("error reading reply channel", zap.Error(err))
		return nil
	}

	correlationHeaders := d.correlationHeaders(message.Headers())

	cmdMsg := commandMessage{sagaID, sagaName, command, correlationHeaders}

	replies, err := handler(ctx, cmdMsg)
	if err != nil {
		logger.Error("saga command handler returned an error", zap.Error(err))
		rerr := d.sendReplies(ctx, replyChannel, []msg.Reply{msg.WithFailure()}, correlationHeaders)
		if rerr != nil {
			logger.Error("error sending replies", zap.Error(rerr))
			return rerr
		}
		return nil
	}

	err = d.sendReplies(ctx, replyChannel, replies, correlationHeaders)
	if err != nil {
		logger.Error("error sending replies", zap.Error(err))
		return err
	}

	return nil
}

func (d *CommandDispatcher) commandMessageInfo(message msg.Message) (string, string, string, error) {
	var err error
	var commandName, sagaID, sagaName string

	commandName, err = message.Headers().GetRequired(msg.MessageCommandName)
	if err != nil {
		d.logger.Error("error reading command name", zap.Error(err))
		return "", "", "", err
	}

	sagaID, err = message.Headers().GetRequired(MessageCommandSagaID)
	if err != nil {
		d.logger.Error("error reading saga id", zap.Error(err))
		return "", "", "", err
	}

	sagaName, err = message.Headers().GetRequired(MessageCommandSagaName)
	if err != nil {
		d.logger.Error("error reading saga name", zap.Error(err))
		return "", "", "", err
	}

	return commandName, sagaID, sagaName, nil
}

func (d *CommandDispatcher) sendReplies(ctx context.Context, replyChannel string, replies []msg.Reply, correlationHeaders msg.Headers) error {
	for _, reply := range replies {
		if err := d.publisher.PublishReply(ctx, reply.Reply(),
			msg.WithHeaders(correlationHeaders),
			msg.WithHeaders(reply.Headers()),
			msg.WithDestinationChannel(replyChannel),
		); err != nil {
			return err
		}
	}

	return nil
}

func (d *CommandDispatcher) correlationHeaders(headers msg.Headers) msg.Headers {
	replyHeaders := make(map[string]string)
	for key, value := range headers {
		if key == msg.MessageCommandName {
			continue
		}

		if strings.HasPrefix(key, msg.MessageCommandPrefix) {
			replyHeader := msg.MessageReplyPrefix + key[len(msg.MessageCommandPrefix):]
			replyHeaders[replyHeader] = value
		}
	}

	return replyHeaders
}
