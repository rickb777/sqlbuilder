package sqlbuilder

import (
	"reflect"
	"testing"
)

type customer struct {
	ID    int
	Name  string
	Phone *string
	Age   int
}

func TestSimpleSelectWithOrder(t *testing.T) {
	c := customer{}

	query, _, dest := MySQL.Select().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		MapSQL("1+1 AS two", nil).
		Order("name", "age").
		Build()

	expectedQuery := "SELECT `id`, `name`, `phone`, `age`, 1+1 AS two FROM `customers` ORDER BY `name`, `age`"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedDest := []interface{}{&c.ID, &c.Name, &c.Phone, &c.Age, &nullDest}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}

func TestSimpleSelectWithLimitOffset(t *testing.T) {
	c := customer{}

	query, _, dest := MySQL.Select().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		Limit(5).
		Offset(10).
		Build()

	expectedQuery := "SELECT `id`, `name`, `phone`, `age` FROM `customers` LIMIT 5 OFFSET 10"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedDest := []interface{}{&c.ID, &c.Name, &c.Phone, &c.Age}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}

func TestSimpleSelectWithJoins(t *testing.T) {
	c := customer{}

	query, _, _ := MySQL.Select().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		Join("INNER JOIN orders ON orders.customer_id = customers.id").
		Join("LEFT JOIN items ON items.order_id = orders.id").
		Build()

	expectedQuery := "SELECT `id`, `name`, `phone`, `age` FROM `customers` INNER JOIN orders ON orders.customer_id = customers.id LEFT JOIN items ON items.order_id = orders.id"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
}

func TestSelectWithWhereMySQL(t *testing.T) {
	c := customer{}

	query, args, _ := MySQL.Select().
		From("customers").As("c").
		Map("c.id", &c.ID).
		Map("c.name", &c.Name).
		Map("c.telephone", &c.Phone).As("phone").
		Map("c.age", &c.Age).
		Where("c.id", "= ?", 9).
		Where("c.name", "IS NOT NULL").
		Where("c.age", "BETWEEN ? AND ?", 10, 20).
		Build()

	expectedQuery := "SELECT `c.id`, `c.name`, `c.telephone` AS `phone`, `c.age` FROM `customers` AS `c` " +
	"WHERE (`c.id` = ?) AND (`c.name` IS NOT NULL) AND (`c.age` BETWEEN ? AND ?)"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9, 10, 20}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestSelectWithGroupMySQL(t *testing.T) {
	var count uint
	query, _, _ := MySQL.Select().From("customers").MapSQL("COUNT(*)", &count).Group("city").Build()
	expectedQuery := "SELECT COUNT(*) FROM `customers` GROUP BY city"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
}

func TestSelectWithWherePostgres(t *testing.T) {
	c := customer{}

	query, args, _ := Postgres.Select().
	From("customers").As("c").
	Map("c.id", &c.ID).
	Map("c.name", &c.Name).
	Map("c.telephone", &c.Phone).As("phone").
	Map("c.age", &c.Age).
	Where("c.id", "= ?", 9).
	Where("c.name", "IS NOT NULL").
	Where("c.age", "BETWEEN ? AND ?", 10, 20).
	Build()

	expectedQuery := `SELECT "c.id", "c.name", "c.telephone" AS "phone", "c.age" FROM "customers" AS "c" `+
	`WHERE ("c.id" = $1) AND ("c.name" IS NOT NULL) AND ("c.age" BETWEEN $2 AND $3)`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9, 10, 20}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}
