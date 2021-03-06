package sqlbuilder

import (
	"strings"
)

// Update returns a new UPDATE statement with the default dialect.
func Update() UpdateStatement {
	return UpdateStatement{dialect: DefaultDialect}
}

type updateSet struct {
	col string
	arg interface{}
	raw bool
}

// UpdateStatement represents an UPDATE statement.
type UpdateStatement struct {
	dialect Dialect
	last    lastWas
	table   name
	sets    []updateSet
	wheres  []where
	args    []interface{}
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s UpdateStatement) Dialect(dialect Dialect) UpdateStatement {
	s.dialect = dialect
	return s
}

// Table returns a new statement with the table to update set to 'table'.
func (s UpdateStatement) Table(table string) UpdateStatement {
	s.table = name{table, ""}
	return s
}

// Set returns a new statement with column 'col' set to value 'val'.
func (s UpdateStatement) Set(col string, val interface{}) UpdateStatement {
	s.sets = append(s.sets, updateSet{col: col, arg: val, raw: false})
	return s
}

// SetSQL returns a new statement with column 'col' set to SQL expression 'sql'.
func (s UpdateStatement) SetSQL(col string, sql string) UpdateStatement {
	s.sets = append(s.sets, updateSet{col: col, arg: sql, raw: true})
	return s
}

// Where returns a new statement with a where-clause consisting of a column, a condition and
// the necessary arguments to that condition.
// For example Where("x", "BETWEEN ? AND ?", 10, 20)
//
// Multiple where-clauses are combined with AND.
func (s UpdateStatement) Where(col, cond string, args ...interface{}) UpdateStatement {
	s.wheres = append(s.wheres, where{col, cond, args})
	return s
}

// WhereEq returns a new statement with condition 'col = ?'. This is a shorthand for Where.
// Multiple where-clauses are combined with AND.
func (s UpdateStatement) WhereEq(col string, args ...interface{}) UpdateStatement {
	return s.Where(col, "=?", args...)
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s UpdateStatement) Build() (query string, args []interface{}) {
	if len(s.sets) == 0 {
		panic("sqlbuilder: UPDATE with no columns set")
	}

	query = "UPDATE " + s.table.String() + " SET "
	var sets []string
	idx := 0

	for _, set := range s.sets {
		var arg string
		if set.raw {
			arg = set.arg.(string)
		} else {
			arg = s.dialect.Placeholder(idx)
			idx++
			args = append(args, set.arg)
		}
		sets = append(sets, set.col+" = "+arg)
	}
	query += strings.Join(sets, ", ")

	if len(s.wheres) > 0 {
		query, args, idx = buildWhereClause(query, args, idx, s.wheres, s.dialect)
	}

	return
}
