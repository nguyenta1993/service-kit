package pgx

import "github.com/gogovan/ggx-kr-service-utils/logger"

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
