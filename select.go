package sqlbuilder

import (
	"strconv"
	"strings"
	"fmt"
)

var nullDest interface{}

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	dbms    DBMS
	last    lastWas
	table   name
	selects []sel
	joins   []join
	wheres  []where
	lock    bool
	limit   *int
	offset  *int
	order   []string
	group   string
	having  string
}

type sel struct {
	col  name
	dest interface{}
	raw  bool
}

type join struct {
	sql  string
	args []interface{}
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
	s.selects = append(s.selects, sel{name{col, ""}, dest, false})
	s.last = lastWasColumnName
	return s
}

// MapSQL is Map without quoting col.
func (s SelectStatement) MapSQL(col string, dest interface{}) SelectStatement {
	if dest == nil {
		dest = nullDest
	}
	s.selects = append(s.selects, sel{name{col, ""}, dest, true})
	s.last = lastWasColumnName
	return s
}

// As modifies the preceding table or column name by setting an alias.
func (s SelectStatement) As(alias string) SelectStatement {
	switch s.last {
	case lastWasTableName:
		s.table = name{s.table.name, alias}
	case lastWasColumnName:
		i := len(s.selects) - 1
		sel := s.selects[i]
		s.selects = s.selects[:i]
		sel.col.alias = alias
		s.selects = append(s.selects, sel)
	}
	s.last = lastWasUnknown
	return s
}

// Join returns a new statement with JOIN expression 'sql'.
func (s SelectStatement) Join(sql string, args ...interface{}) SelectStatement {
	s.joins = append(s.joins, join{sql, args})
	return s
}

// Where returns a new statement with condition on a specified column.
// Multiple conditions are combined with AND.
func (s SelectStatement) Where(col, cond string, args ...interface{}) SelectStatement {
	s.wheres = append(s.wheres, where{col, cond, args})
	return s
}

// WhereEq returns a new statement with condition 'col = ?'.
// Multiple conditions are combined with AND.
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

// Order returns a new statement with ordering 'order', which may be a list of column names.
// Only the last Order() is used.
func (s SelectStatement) Order(order ...string) SelectStatement {
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

	if len(s.selects) > 0 {
		for _, sel := range s.selects {
			if sel.raw {
				cols = append(cols, sel.col.String())
			} else {
				cols = append(cols, sel.col.Quoted(s.dbms.Dialect))
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

	query = fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(cols, ", "),
		s.table.Quoted(s.dbms.Dialect))

	for _, join := range s.joins {
		sql := join.sql
		for _, arg := range join.args {
			sql = strings.Replace(sql, "?", s.dbms.Dialect.Placeholder(idx), 1)
			idx++
			args = append(args, arg)
		}
		query += " " + sql
	}

	if len(s.wheres) > 0 {
		query, args, idx = buildWhereClause(query, args, idx, s.wheres, s.dbms.Dialect)
	}

	if len(s.order) > 0 {
		quoted := quoteStrings(s.order, s.dbms.Dialect)
		query += " ORDER BY " + strings.Join(quoted, ", ")
	}

	if s.group != "" {
		query += " GROUP BY " + s.group
	}

	if s.having != "" {
		query += " HAVING " + s.having
	}

	if s.limit != nil {
		query += " LIMIT " + strconv.Itoa(*s.limit)
	}

	if s.offset != nil {
		query += " OFFSET " + strconv.Itoa(*s.offset)
	}

	if s.lock {
		query += " FOR UPDATE"
	}

	return
}
