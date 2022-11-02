package database

import (
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"

	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/jmoiron/sqlx"
)

// MustConnect open a database with current context and support tracing
func MustConnect(dbType, connectionString string) *sqlx.DB {
	return otelsqlx.MustConnect(dbType, connectionString, otelsql.WithAttributes(semconv.DBSystemKey.String(dbType)))
}
