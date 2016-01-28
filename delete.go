package sqlbuilder

type deleteSet struct {
	col string
	arg interface{}
	raw bool
}

// DeleteStatement represents an DELETE statement.
type DeleteStatement struct {
	dbms   DBMS
	table  string
	wheres []where
	args   []interface{}
}

// From returns a new statement with the table to delete set to 'table'.
func (s DeleteStatement) From(table string) DeleteStatement {
	s.table = table
	return s
}

// Where returns a new statement with condition 'cond'.
// Multiple Where() are combined with AND.
// Be careful to use this always; a delete without a where clause is probably incorrect.
func (s DeleteStatement) Where(col, cond string, args ...interface{}) DeleteStatement {
	s.wheres = append(s.wheres, where{col, cond, args})
	return s
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s DeleteStatement) Build() (query string, args []interface{}) {
	if len(s.wheres) == 0 {
		panic("sqlbuilder: DELETE with no where clauses")
	}

	query = "DELETE FROM " + s.dbms.Dialect.Quote(s.table)

	query, args, _ = buildWhereClause(query, args, 0, s.wheres, s.dbms.Dialect)

	return
}
