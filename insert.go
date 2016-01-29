package sqlbuilder

import (
	"fmt"
	"strings"
)

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
	dbms  DBMS
	last  lastWas
	table name
	sets  []insertSet
	rets  []insertRet
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
		//case lastWasColumnName:
		//	i := len(s.selects) - 1
		//	sel := s.selects[i]
		//	s.selects = s.selects[:i]
		//	sel.col.alias = alias
		//	s.selects = append(s.selects, sel)
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
	s.rets = append(s.rets, insertRet{sql: s.dbms.Dialect.Quote(col), dest: dest})
	return s
}

// ReturnSQL is Return without quoting the argument.
func (s InsertStatement) ReturnSQL(sql string, dest interface{}) InsertStatement {
	s.rets = append(s.rets, insertRet{sql: sql, dest: dest})
	return s
}

// Build builds the SQL query. It returns the SQL query and the argument slice.
func (s InsertStatement) Build() (query string, args []interface{}, dest []interface{}) {
	var cols, vals []string
	idx := 0

	for _, set := range s.sets {
		cols = append(cols, s.dbms.Dialect.Quote(set.col))

		if set.raw {
			vals = append(vals, set.arg.(string))
		} else {
			args = append(args, set.arg)
			vals = append(vals, s.dbms.Dialect.Placeholder(idx))
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
		s.table.Quoted(s.dbms.Dialect),
		strings.Join(cols, ", "),
		strings.Join(vals, ", "),
		returning)

	return
}
