package saga

import (
	"context"

	"github.com/tikivn/s14e-backend-utils/saga/core"
)

// RemoteStepActionOption options for remoteStepAction
type RemoteStepActionOption func(action *remoteStepAction)

// WithRemoteStepPredicate sets a predicate function for the action
func WithRemoteStepPredicate(predicate func(context.Context, core.SagaData) bool) RemoteStepActionOption {
	return func(step *remoteStepAction) {
		step.predicate = predicate
	}
}
