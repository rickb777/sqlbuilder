package sqlbuilder

import (
	"fmt"
	"strings"
)

// Insert returns a new INSERT statement with the default dialect.
func Insert() InsertStatement {
	return InsertStatement{dialect: DefaultDialect}
}

type insertSet struct {
	col string
	arg interface{}
	raw bool
}

type insertRet struct {
	sql  string
	dest interface{}
}

// InsertStatement represents an INSERT statement.
type InsertStatement struct {
	dialect Dialect
	last    lastWas
	table   name
	sets    []insertSet
	rets    []insertRet
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s InsertStatement) Dialect(dialect Dialect) InsertStatement {
	s.dialect = dialect
	return s
}

// Into returns a new statement with the table to insert into set to 'table'.
func (s InsertStatement) Into(table string) InsertStatement {
	s.table = name{table, ""}
	s.last = lastWasTableName
	return s
}

// As modifies the preceding table or column name by setting an alias.
func (s InsertStatement) As(alias string) InsertStatement {
	switch s.last {
	case lastWasTableName:
		s.table = name{s.table.name, alias}
	}
	s.last = lastWasUnknown
	return s
}

// Set returns a new statement with column 'col' set to value 'val'.
func (s InsertStatement) Set(col string, val interface{}) InsertStatement {
	s.sets = append(s.sets, insertSet{col, val, false})
	return s
}

// SetSQL returns a new statement with column 'col' set to the raw SQL expression 'sql'.
func (s InsertStatement) SetSQL(col, sql string) InsertStatement {
	s.sets = append(s.sets, insertSet{col, sql, true})
	return s
}

// Return returns a new statement with a RETURNING clause.
func (s InsertStatement) Return(col string, dest interface{}) InsertStatement {
	s.rets = append(s.rets, insertRet{sql: col, dest: dest})
	return s
}

// Build builds the SQL query. It returns the SQL query and the argument slice.
func (s InsertStatement) Build() (query string, args []interface{}, dest []interface{}) {
	var cols, vals []string
	idx := 0

	for _, set := range s.sets {
		cols = append(cols, set.col)

		if set.raw {
			vals = append(vals, set.arg.(string))
		} else {
			args = append(args, set.arg)
			vals = append(vals, s.dialect.Placeholder(idx))
			idx++
		}
	}

	returning := ""
	if len(s.rets) > 0 {
		var args []string
		for _, ret := range s.rets {
			args = append(args, ret.sql)
			dest = append(dest, ret.dest)
		}
		returning = " RETURNING " + strings.Join(args, ", ")
	}

	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)%s",
		s.table,
		strings.Join(cols, ", "),
		strings.Join(vals, ", "),
		returning)

	return
}
