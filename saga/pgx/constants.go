package pgx

import (
	"errors"
)

type contextKey int

const (
	DefaultSagaInstanceTableName = "saga_instances"

	CreateSagaInstancesTableSQL = `CREATE TABLE %s (
    saga_name      text        NOT NULL,
    saga_id        text        NOT NULL,
    saga_data_name text        NOT NULL,
    saga_data      bytea       NOT NULL,
    current_step   int         NOT NULL,
    end_state      boolean     NOT NULL,
    compensating   boolean     NOT NULL,
    modified_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (saga_name, saga_id)
)`

	findSagaInstanceSQL   = "SELECT saga_data_name, saga_data, current_step, end_state, compensating FROM %s WHERE saga_name = $1 AND saga_id = $2 LIMIT 1"
	saveSagaInstanceSQL   = "INSERT INTO %s (saga_name, saga_id, saga_data_name, saga_data, current_step, end_state, compensating, modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)"
	updateSagaInstanceSQL = "UPDATE %s SET saga_data = $1, current_step = $2, end_state = $3, compensating = $4, modified_at = CURRENT_TIMESTAMP WHERE saga_name = $5 AND saga_id = $6"

	pgxTxKey = contextKey(5432)
)

var ErrTxNotInContext = errors.New("pgx.Tx is not set for session")
var ErrInvalidTxValue = errors.New("tx value is not a pgx.Tx type")
