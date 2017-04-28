package sqlbuilder

import (
	"strconv"
)

// Dialect represents a SQL dialect.
type Dialect interface {
	// Placeholder returns the placeholder binding string for parameter at index idx.
	Placeholder(idx int) string
}

type MySQLDialect struct{}
type PostgresDialect struct{}

var (
	MySQL    MySQLDialect    // MySQL
	SQLite   MySQLDialect    // SQLite (same as MySQL)
	Postgres PostgresDialect // Postgres
)

var DefaultDialect = MySQL // Default dialect

func (dialect MySQLDialect) Placeholder(idx int) string {
	return "?"
}

func (dialect PostgresDialect) Placeholder(idx int) string {
	return "$" + strconv.Itoa(idx+1)
}

