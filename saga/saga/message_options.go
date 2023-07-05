package saga

import "github.com/nguyenta1993/service-kit/saga/msg"

// WithSagaInfo is an option to set additional Saga specific headers
func WithSagaInfo(instance *Instance) msg.MessageOption {
	return msg.WithHeaders(map[string]string{
		MessageCommandSagaID:   instance.sagaID,
		MessageCommandSagaName: instance.sagaName,
	})
}
