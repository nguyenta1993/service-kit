package saga

import (
	"context"

	"github.com/nguyenta1993/service-kit/saga/core"
)

// RemoteStepActionOption options for remoteStepAction
type RemoteStepActionOption func(action *remoteStepAction)

// WithRemoteStepPredicate sets a predicate function for the action
func WithRemoteStepPredicate(predicate func(context.Context, core.SagaData) bool) RemoteStepActionOption {
	return func(step *remoteStepAction) {
		step.predicate = predicate
	}
}
