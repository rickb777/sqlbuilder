package sqlbuilder

import (
	"fmt"
	"strings"
)

type join struct {
	op       string
	table    name
	onL, onR name
	using    []string
	dialect  Dialect
}

// Natural precedes Join when required. Any of the other modifiers Left, LeftOuter,
// Right, RightOuter, FullOuter, Inner, Cross may also be used.
func (s SelectStatement) Natural() SelectStatement {
	s.joinNat = "NATURAL "
	return s
}

// Left precedes Join when required. Only one join modifier can be used.
func (s SelectStatement) Left() SelectStatement {
	s.joinOp = "LEFT "
	return s
}

// LeftOuter precedes Join when required. Only one join modifier can be used.
func (s SelectStatement) LeftOuter() SelectStatement {
	s.joinOp = "LEFT OUTER "
	return s
}

// Right precedes Join when required. Only one join modifier can be used.
// SQLite does not support right join.
func (s SelectStatement) Right() SelectStatement {
	s.joinOp = "RIGHT "
	return s
}

// RightOuter precedes Join when required. Only one join modifier can be used.
// SQLite does not support right outer join.
func (s SelectStatement) RightOuter() SelectStatement {
	s.joinOp = "RIGHT OUTER "
	return s
}

// FullOuter precedes Join when required. Only one join modifier can be used.
// SQLite does not support full outer join.
func (s SelectStatement) FullOuter() SelectStatement {
	s.joinOp = "FULL OUTER "
	return s
}

// Inner precedes Join when required. Only one join modifier can be used.
func (s SelectStatement) Inner() SelectStatement {
	s.joinOp = "INNER "
	return s
}

// Cross precedes Join when required. Only one join modifier can be used and this is
// not compatible with natural join.
func (s SelectStatement) Cross() SelectStatement {
	s.joinNat = ""
	s.joinOp = "CROSS "
	return s
}

// Join sets the table name for the current join.
func (s SelectStatement) Join(table string) SelectStatement {
	s.joinTbl = name{table, ""}
	s.last = lastWasJoinTableName
	return s
}

// On completes a JOIN clause with the necessary constraint.
// When required, another join can immediately follow this.
func (s SelectStatement) On(onL, onR string) SelectStatement {
	op := s.joinNat + s.joinOp + "JOIN"
	j := join{op, s.joinTbl, splitAsName(onL), splitAsName(onR), nil, s.dialect}
	s.joins = append(s.joins, j)
	s.joinNat = ""
	s.joinOp = ""
	s.joinTbl = name{}
	return s
}

// Using completes a JOIN clause with the necessary columns.
// When required, another join can immediately follow this.
func (s SelectStatement) Using(col ...string) SelectStatement {
	op := s.joinNat + s.joinOp + "JOIN"
	j := join{op, s.joinTbl, name{}, name{}, col, s.dialect}
	s.joins = append(s.joins, j)
	s.joinNat = ""
	s.joinOp = ""
	s.joinTbl = name{}
	return s
}

func (j join) String() string {
	tbl := j.table.QuotedAs(j.dialect)
	if len(j.using) > 0 {
		cols := strings.Join(j.using, ", ")
		return fmt.Sprintf("\n %s %s USING (%s)", j.op, tbl, cols)
	} else {
		onL := j.onL.QuotedDot(j.dialect)
		onR := j.onR.QuotedDot(j.dialect)
		return fmt.Sprintf("\n %s %s ON %s = %s", j.op, tbl, onL, onR)
	}
}
