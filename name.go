package sqlbuilder

import "strings"

type lastWas int

const (
	lastWasUnknown lastWas = iota
	lastWasTableName
	lastWasJoinTableName
	lastWasColumnName
)

type name struct {
	name, alias string
}

func splitAsName(s string) name {
	a := strings.Split(s, ".")
	switch len(a) {
	case 1:
		return name{s, ""}
	case 2:
		return name{a[0], a[1]}
	}
	panic("No way to split '" + s + "' at the dot")
}

func (n name) QuotedAs(dialect Dialect) string {
	return n.quotedX(dialect, " AS ")
}

func (n name) QuotedDot(dialect Dialect) string {
	return n.quotedX(dialect, ".")
}

func (n name) quotedX(dialect Dialect, sep string) string {
	if n.alias == "" {
		return n.name
	}
	return n.name + sep + n.alias
}

func (n name) String() string {
	qn := n.name
	if n.alias == "" {
		return qn
	}
	return qn + " AS " + n.alias
}
