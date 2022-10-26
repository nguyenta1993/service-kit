package saga

import (
	"github.com/tikivn/s14e-backend-utils/saga/core"
	"github.com/tikivn/s14e-backend-utils/saga/msg"
)

// Command is a core.Command with message header information
type Command interface {
	SagaID() string
	SagaName() string
	Command() core.Command
	Headers() msg.Headers
}

type commandMessage struct {
	sagaID   string
	sagaName string
	command  core.Command
	headers  msg.Headers
}

func (m commandMessage) SagaID() string {
	return m.sagaID
}

func (m commandMessage) SagaName() string {
	return m.sagaName
}

func (m commandMessage) Command() core.Command {
	return m.command
}

func (m commandMessage) Headers() msg.Headers {
	return m.headers
}
