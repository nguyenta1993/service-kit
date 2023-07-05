package pgx

import "github.com/nguyenta1993/service-kit/logger"

type SagaInstanceStoreOption func(*SagaInstanceStore)

func WithSagaInstanceStoreTableName(tableName string) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.tableName = tableName
	}
}

func WithSagaInstanceStoreLogger(logger logger.Logger) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.logger = logger
	}
}
