package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

var nullDest interface{}

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	dbms    DBMS
	last    lastWas
	table   name
	selects []sel
	joinNat string
	joinOp  string
	joinTbl name
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
	op       string
	table    name
	onL, onR name
	using    []string
	dialect  Dialect
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
	case lastWasJoinTableName:
		s.joinTbl = name{s.joinTbl.name, alias}
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

// Natural precedes Join when required.
func (s SelectStatement) Natural() SelectStatement {
	s.joinNat = "NATURAL "
	return s
}

// Left precedes Join when required.
func (s SelectStatement) Left() SelectStatement {
	s.joinOp = "LEFT "
	return s
}

// LeftOuter precedes Join when required.
func (s SelectStatement) LeftOuter() SelectStatement {
	s.joinOp = "LEFT OUTER "
	return s
}

// Inner precedes Join when required.
func (s SelectStatement) Inner() SelectStatement {
	s.joinOp = "INNER "
	return s
}

// Cross precedes Join when required.
func (s SelectStatement) Cross() SelectStatement {
	s.joinOp = "CROSS "
	return s
}

// Join sets the table name for the current join.
func (s SelectStatement) Join(table string) SelectStatement {
	s.joinTbl = name{table, ""}
	s.last = lastWasJoinTableName
	return s
}

// On completes a JOIN clause with the necessary constraint
func (s SelectStatement) On(onL, onR string) SelectStatement {
	op := s.joinNat + s.joinOp + "JOIN"
	j := join{op, s.joinTbl, splitAsName(onL), splitAsName(onR), nil, s.dbms.Dialect}
	s.joins = append(s.joins, j)
	s.joinNat = ""
	s.joinOp = ""
	s.joinTbl = name{}
	return s
}

// Using completes a JOIN clause with the necessary columns
func (s SelectStatement) Using(col ...string) SelectStatement {
	op := s.joinNat + s.joinOp + "JOIN"
	j := join{op, s.joinTbl, name{}, name{}, col, s.dbms.Dialect}
	s.joins = append(s.joins, j)
	s.joinNat = ""
	s.joinOp = ""
	s.joinTbl = name{}
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

	if len(s.selects) > 0 {
		for _, sel := range s.selects {
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

	query = fmt.Sprintf("SELECT %s\n FROM %s",
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

func (j join) String() string {
	tbl := j.table.QuotedAs(j.dialect)
	if len(j.using) > 0 {
		cols := strings.Join(quoteStrings(j.using, j.dialect), ", ")
		return fmt.Sprintf("\n %s %s USING %s", j.op, tbl, cols)
	} else {
		onL := j.onL.QuotedDot(j.dialect)
		onR := j.onR.QuotedDot(j.dialect)
		return fmt.Sprintf("\n %s %s ON %s = %s", j.op, tbl, onL, onR)
	}
}
