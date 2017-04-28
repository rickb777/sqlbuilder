sqlbuilder
==========

[![Travis CI status](https://api.travis-ci.org/rickb777/sqlbuilder.svg)](https://travis-ci.org/rickb777/sqlbuilder)

`sqlbuilder` is a Go library for constructing SQL queries using a fluent API.

The master branch tracks version 3. The latest stable version is
[2.0.0](https://github.com/thcyron/sqlbuilder/tree/v2.0.0/).

`sqlbuilder` follows [Semantic Versioning](http://semver.org/).

Installation
------------

    go get github.com/rickb777/sqlbuilder

Examples
--------

**SELECT**

```go
query, args, dest := sqlbuilder.Select().
        From("customers").
        Map("id", &customer.ID).
        Map("name", &customer.Name).
        Map("telephone", &customer.Phone).As("phone").
        OrderBy("id DESC").
        Limit(1).
        Build()
err := db.QueryRow(query, args...).Scan(dest...)
```

Joins have a fluent style:

```go
query, args, dest := sqlbuilder.Select().
        From("customers").As("c").
        Inner().Join("orders").As("o").On("o.customer_id", "c.id").
        OrderBy("o.total_price").
        Map("c.id", &cview.ID).
        Map("c.name", &cview.Name).
        Map("c.telephone", &cview.Phone).
        Map("o.total_price", &cview.TotalPrice).
        Build()
```

**INSERT**

```go
query, args := sqlbuilder.Insert().
        Into("customers").
        Set("name", "John").
        Set("phone", "555").
        Build()
err := db.Exec(query, args...)
```

**UPDATE**

```go
query, args := sqlbuilder.Update().
        Table("customers").
        Set("name", "John").
        Set("phone", "555").
        Where("id", "= ?", 1).
        Build()
err := db.Exec(query, args...)
```

**DELETE**

```go
query, args := sqlbuilder.Delete().
        From("customers").
        WhereEq("id", 1).
        Build()
err := db.Exec(query, args...)
```

Supported DBMS
--------------

`sqlbuilder` supports building queries for MySQL, SQLite, and Postgres databases. You
can set the default dialect with:

```go
sqlbuilder.DefaultDialect = sqlbuilder.Postgres
sqlbuilder.Select().From("...")...
```

Or you can specify the dialect explicitly:

```go
sqlbuilder.Select().Dialect(sqlbuilder.Postgres).From("...")...
```

Documentation
-------------

Documentation is available at [Godoc](http://godoc.org/github.com/rickb777/sqlbuilder).

Licence
-------

`sqlbuilder` is licensed under the MIT License.
