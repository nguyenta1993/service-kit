package saga

import (
	"context"

	"github.com/nguyenta1993/service-kit/saga/core"
)

// Step interface for local, remote, ...other saga steps
type Step interface {
	hasInvocableAction(ctx context.Context, sagaData core.SagaData, compensating bool) bool
	getReplyHandler(replyName string, compensating bool) func(ctx context.Context, data core.SagaData, reply core.Reply) error
	execute(ctx context.Context, sagaData core.SagaData, compensating bool) func(results *stepResults)
}
