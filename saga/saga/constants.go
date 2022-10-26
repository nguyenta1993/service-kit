package saga

import "github.com/tikivn/s14e-backend-utils/saga/msg"

const (
	notCompensating = false
	isCompensating  = true
)

// LifecycleHook type for hooking in custom code at specific stages of a saga
type LifecycleHook int

// Definition lifecycle hooks
const (
	SagaStarting LifecycleHook = iota
	SagaCompleted
	SagaCompensated
)

// Saga message headers
const (
	MessageCommandSagaID   = msg.MessageCommandPrefix + "SAGA_ID"
	MessageCommandSagaName = msg.MessageCommandPrefix + "SAGA_NAME"
	MessageCommandResource = msg.MessageCommandPrefix + "RESOURCE"

	MessageReplySagaID   = msg.MessageReplyPrefix + "SAGA_ID"
	MessageReplySagaName = msg.MessageReplyPrefix + "SAGA_NAME"
)
