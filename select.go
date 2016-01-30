package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

var nullDest interface{}

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	dbms     DBMS
	distinct string
	last     lastWas
	table    name
	columns  []column
	joinNat  string
	joinOp   string
	joinTbl  name
	joins    []join
	wheres   []where
	lock     bool
	limit    *int
	offset   *int
	order    []string
	group    string
	having   string
}

type column struct {
	col  name
	dest interface{}
	raw  bool
}

// Distinct modifies the select to remove duplicates from the results.
func (s SelectStatement) Distinct() SelectStatement {
	s.distinct = "DISTINCT "
	return s
}

// From returns a new statement with the table to select from set to 'table'.
func (s SelectStatement) From(table string) SelectStatement {
	s.table = name{table, ""}
	s.last = lastWasTableName
	return s
}

// Map returns a new statement with column 'col' selected and scanned
// into 'dest'. 'dest' may be nil if the value should not be scanned.
func (s SelectStatement) Map(col string, dest interface{}) SelectStatement {
	if dest == nil {
		dest = nullDest
	}
	s.columns = append(s.columns, column{name{col, ""}, dest, false})
	s.last = lastWasColumnName
	return s
}

// MapSQL is Map without quoting col.
func (s SelectStatement) MapSQL(col string, dest interface{}) SelectStatement {
	if dest == nil {
		dest = nullDest
	}
	s.columns = append(s.columns, column{name{col, ""}, dest, true})
	s.last = lastWasColumnName
	return s
}

// As modifies the preceding table or column name by setting an alias.
func (s SelectStatement) As(alias string) SelectStatement {
	switch s.last {
	case lastWasTableName:
		s.table = name{s.table.name, alias}
	case lastWasJoinTableName:
		s.joinTbl = name{s.joinTbl.name, alias}
	case lastWasColumnName:
		i := len(s.columns) - 1
		sel := s.columns[i]
		s.columns = s.columns[:i]
		sel.col.alias = alias
		s.columns = append(s.columns, sel)
	}
	s.last = lastWasUnknown
	return s
}

// Where returns a new statement with a where-clause consisting of a column, a condition and
// the necessary arguments to that condition.
// For example Where("x", "BETWEEN ? AND ?", 10, 20)
//
// Multiple where-clauses are combined with AND.
func (s SelectStatement) Where(col, cond string, args ...interface{}) SelectStatement {
	s.wheres = append(s.wheres, where{col, cond, args})
	return s
}

// WhereEq returns a new statement with condition 'col = ?'. This is a shorthand for Where.
// Multiple where-clauses are combined with AND.
func (s SelectStatement) WhereEq(col string, args ...interface{}) SelectStatement {
	return s.Where(col, "=?", args...)
}

// Limit returns a new statement with the limit set to 'limit'.
func (s SelectStatement) Limit(limit int) SelectStatement {
	s.limit = &limit
	return s
}

// Offset returns a new statement with the offset set to 'offset'.
func (s SelectStatement) Offset(offset int) SelectStatement {
	s.offset = &offset
	return s
}

// OrderBy returns a new statement with ordering 'order', which may be a list of column names.
// Only the last OrderBy() is used.
func (s SelectStatement) OrderBy(order ...string) SelectStatement {
	s.order = order
	return s
}

// Group returns a new statement with grouping 'group'.
// Only the last Group() is used.
func (s SelectStatement) Group(group string) SelectStatement {
	s.group = group
	return s
}

// Having returns a new statement with HAVING condition 'having'.
// Only the last Having() is used.
func (s SelectStatement) Having(having string) SelectStatement {
	s.having = having
	return s
}

// Lock returns a new statement with FOR UPDATE locking.
func (s SelectStatement) Lock() SelectStatement {
	s.lock = true
	return s
}

// Build builds the SQL query. It returns the query, the argument slice,
// and the destination slice.
func (s SelectStatement) Build() (query string, args []interface{}, dest []interface{}) {
	var cols []string
	idx := 0

	if len(s.columns) > 0 {
		for _, sel := range s.columns {
			if sel.raw {
				cols = append(cols, sel.col.String())
			} else {
				cols = append(cols, sel.col.QuotedAs(s.dbms.Dialect))
			}
			if sel.dest == nil {
				dest = append(dest, &nullDest)
			} else {
				dest = append(dest, sel.dest)
			}
		}
	} else {
		cols = append(cols, "1")
		dest = append(dest, &nullDest)
	}

	query = fmt.Sprintf("SELECT %s%s\n FROM %s",
		s.distinct,
		strings.Join(cols, ", "),
		s.table.QuotedAs(s.dbms.Dialect))

	for _, join := range s.joins {
		query += join.String()
	}

	if len(s.wheres) > 0 {
		query, args, idx = buildWhereClause(query, args, idx, s.wheres, s.dbms.Dialect)
	}

	if len(s.order) > 0 {
		quoted := quoteStrings(s.order, s.dbms.Dialect)
		query += "\n ORDER BY " + strings.Join(quoted, ", ")
	}

	if s.group != "" {
		query += "\n GROUP BY " + s.group
	}

	if s.having != "" {
		query += "\n HAVING " + s.having
	}

	if s.limit != nil {
		query += "\n LIMIT " + strconv.Itoa(*s.limit)
	}

	if s.offset != nil {
		query += "\n OFFSET " + strconv.Itoa(*s.offset)
	}

	if s.lock {
		query += "\n FOR UPDATE"
	}

	return
}
