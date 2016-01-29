package sqlbuilder

type lastWas int

const (
	lastWasUnknown lastWas = iota
	lastWasTableName
	lastWasColumnName
)

type name struct {
	name, alias string
}

func (n name) Quoted(dialect Dialect) string {
	qn := dialect.Quote(n.name)
	if n.alias == "" {
		return qn
	}
	return qn + " AS " + dialect.Quote(n.alias)
}

func (n name) String() string {
	qn := n.name
	if n.alias == "" {
		return qn
	}
	return qn + " AS " + n.alias
}
