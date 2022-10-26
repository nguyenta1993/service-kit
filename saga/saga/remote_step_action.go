package saga

import (
	"context"

	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/msg"
)

type remoteStepAction struct {
	predicate func(context.Context, core.SagaData) bool
	handler   func(context.Context, core.SagaData) msg.DomainCommand
}

func (a *remoteStepAction) isInvocable(ctx context.Context, sagaData core.SagaData) bool {
	if a.predicate == nil {
		return true
	}

	return a.predicate(ctx, sagaData)
}

func (a *remoteStepAction) execute(ctx context.Context, sagaData core.SagaData) msg.DomainCommand {
	return a.handler(ctx, sagaData)
}
