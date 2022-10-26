package msg

import (
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/core"
)

// RegisterTypes should be called after registering a new marshaller; especially after registering a new default
func RegisterTypes() {
	// Need to register the success and failure messages with the msgpack marshaller
	core.RegisterReplies(Success{}, Failure{})
}
