package saga

import (
	"github.com/tikivn/s14e-backend-utils/saga/core"
	"github.com/tikivn/s14e-backend-utils/saga/msg"
)

type stepResults struct {
	commands           []msg.DomainCommand
	updatedSagaData    core.SagaData
	updatedStepContext stepContext
	local              bool
	failure            error
}
