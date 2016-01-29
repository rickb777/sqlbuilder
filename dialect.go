package sqlbuilder

import (
	"strconv"
	"strings"
)

const mysqlQuote = "`"
const postgresQuote = `"`

type Dialect interface {
	// Placeholder returns the placeholder string for the given index.
	Placeholder(idx int) string

	// Quote returns s quoted.
	Quote(s string) string
}

type MySQLDialect struct{}
type PostgresDialect struct{}

func (dialect MySQLDialect) Placeholder(idx int) string {
	return "?"
}

func (dialect PostgresDialect) Placeholder(idx int) string {
	return "$" + strconv.Itoa(idx+1)
}

func quote(s, quote string) string {
	return quote + strings.Replace(s, quote, "\\"+quote, -1) + quote
}

func (dialect MySQLDialect) Quote(s string) string {
	return quote(s, mysqlQuote)
}

func (dialect PostgresDialect) Quote(s string) string {
	return quote(s, postgresQuote)
}
