package saga

import (
	"github.com/nguyenta1993/service-kit/saga/core"
	"github.com/nguyenta1993/service-kit/saga/msg"
)

type stepResults struct {
	commands           []msg.DomainCommand
	updatedSagaData    core.SagaData
	updatedStepContext stepContext
	local              bool
	failure            error
}
