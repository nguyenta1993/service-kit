package msg

import (
	"github.com/nguyenta1993/service-kit/saga/core"
)

// RegisterTypes should be called after registering a new marshaller; especially after registering a new default
func RegisterTypes() {
	// Need to register the success and failure messages with the msgpack marshaller
	core.RegisterReplies(Success{}, Failure{})
}
