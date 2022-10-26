package saga

import (
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/msg"
)

type stepResults struct {
	commands           []msg.DomainCommand
	updatedSagaData    core.SagaData
	updatedStepContext stepContext
	local              bool
	failure            error
}
