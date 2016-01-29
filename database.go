package sqlbuilder

// DBMS represents a DBMS.
type DBMS struct {
	Dialect Dialect
}

var (
	// The simplest syntax variant, compatible with MySQL and SQLite.
	// Identifiers must be simple alphanumeric strings (with underscore optional).
	MySQL = DBMS{PlainDialect{}}
	// The syntax variant needed by MySQL with identifiers enclosed in back-ticks.
	MySQLQuoted = DBMS{MySQLDialect{}}
	// The syntax variant needed by PostgreSQL
	Postgres = DBMS{PostgresDialect{}}
)

// Select returns a new SELECT statement.
func (dbms DBMS) Select() SelectStatement {
	return SelectStatement{dbms: dbms}
}

// Insert returns a new INSERT statement.
func (dbms DBMS) Insert() InsertStatement {
	return InsertStatement{dbms: dbms}
}

// Update returns a new UPDATE statement.
func (dbms DBMS) Update() UpdateStatement {
	return UpdateStatement{dbms: dbms}
}

// Delete returns a new UPDATE statement.
func (dbms DBMS) Delete() DeleteStatement {
	return DeleteStatement{dbms: dbms}
}

// DefaultDBMS is the DBMS used by the package-level Select,
// Insert, Update and Delete functions.
var DefaultDBMS = MySQL

// Select returns a new SELECT statement using the default DBMS.
func Select() SelectStatement {
	return DefaultDBMS.Select()
}

// Insert returns a new INSERT statement using the default DBMS.
func Insert() InsertStatement {
	return DefaultDBMS.Insert()
}

// Update returns a new UPDATE statement using the default DBMS.
func Update() UpdateStatement {
	return DefaultDBMS.Update()
}

// Delete returns a new DELETE statement using the default DBMS.
func Delete() DeleteStatement {
	return DefaultDBMS.Delete()
}
