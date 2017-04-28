package sqlbuilder

type deleteSet struct {
	col string
	arg interface{}
	raw bool
}

// Delete returns a new DELETE statement with the default dialect.
func Delete() DeleteStatement {
	return DeleteStatement{dialect: DefaultDialect}
}

// DeleteStatement represents an DELETE statement.
type DeleteStatement struct {
	dialect Dialect
	last    lastWas
	table   name
	wheres  []where
	args    []interface{}
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s DeleteStatement) Dialect(dialect Dialect) DeleteStatement {
	s.dialect = dialect
	return s
}

// From returns a new statement with the table to delete set to 'table'.
func (s DeleteStatement) From(table string) DeleteStatement {
	s.table = name{table, ""}
	s.last = lastWasTableName
	return s
}

// As modifies the preceding table or column name by setting an alias.
func (s DeleteStatement) As(alias string) DeleteStatement {
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

// Where returns a new statement with a where-clause consisting of a column, a condition and
// the necessary arguments to that condition.
// For example Where("x", "BETWEEN ? AND ?", 10, 20)
//
// Multiple where-clauses are combined with AND.
// Be careful to use this always; a delete without a where clause is probably incorrect.
func (s DeleteStatement) Where(col, cond string, args ...interface{}) DeleteStatement {
	s.wheres = append(s.wheres, where{col, cond, args})
	return s
}

// WhereEq returns a new statement with condition 'col = ?'. This is a shorthand for Where.
// Multiple where-clauses are combined with AND.
func (s DeleteStatement) WhereEq(col string, args ...interface{}) DeleteStatement {
	return s.Where(col, "=?", args...)
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s DeleteStatement) Build() (query string, args []interface{}) {
	if len(s.wheres) == 0 {
		panic("sqlbuilder: DELETE with no where clauses")
	}

	query = "DELETE FROM " + s.table.String()

	query, args, _ = buildWhereClause(query, args, 0, s.wheres, s.dialect)

	return
}
