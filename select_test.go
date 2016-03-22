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

func TestSimpleSelectWithOrderAndLock(t *testing.T) {
	c := customer{}

	query, _, dest := MySQLQuoted.Select().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		Map("1+1", nil).As("two").
		OrderBy("name", "age").
		Lock().
		Build()

	expectedQuery := "SELECT `id`, `name`, `phone`, `age`, `1+1` AS `two`\n FROM `customers`\n" +
		" ORDER BY `name`, `age`\n FOR UPDATE"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedDest := []interface{}{&c.ID, &c.Name, &c.Phone, &c.Age, &nullDest}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}

func TestSimpleSelectDistinctWithLimitOffset(t *testing.T) {
	c := customer{}

	query, _, dest := MySQL.Select().Distinct().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		Limit(5).
		Offset(10).
		Build()

	expectedQuery := "SELECT DISTINCT id, name, phone, age\n FROM customers\n LIMIT 5\n OFFSET 10"
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
		From("customers").As("c").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("age", &c.Age).
		Inner().Join("orders").As("o").On("o.customer_id", "c.id").
		Left().Join("items").Using("id").
		Build()

	expectedQuery := "SELECT id, name, phone, age\n" +
		" FROM customers AS c\n" +
		" INNER JOIN orders AS o ON o.customer_id = c.id\n" +
		" LEFT JOIN items USING (id)"
	if query != expectedQuery {
		t.Errorf("bad query: |%s|", query)
	}
}

func TestSelectWithWhereMySQL(t *testing.T) {
	c := customer{}

	query, args, _ := MySQLQuoted.Select().
		From("customers").As("c").
		Map("c.id", &c.ID).
		Map("c.name", &c.Name).
		Map("c.telephone", &c.Phone).As("phone").
		Map("c.age", &c.Age).
		Cross().Join("orders").As("o").On("o.customer_id", "c.id").
		Where("c.id", "= ?", 9).
		Where("c.name", "IS NOT NULL").
		Where("c.age", "BETWEEN ? AND ?", 10, 20).
		Build()

	expectedQuery := "SELECT `c.id`, `c.name`, `c.telephone` AS `phone`, `c.age`\n" +
		" FROM `customers` AS `c`\n" +
		" CROSS JOIN `orders` AS `o` ON `o`.`customer_id` = `c`.`id`\n" +
		" WHERE (`c.id` = ?) AND (`c.name` IS NOT NULL) AND (`c.age` BETWEEN ? AND ?)"
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
	expectedQuery := "SELECT COUNT(*)\n FROM customers\n GROUP BY city"
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
		Cross().Join("orders").As("o").On("o.customer_id", "c.id").
		Where("c.id", "= ?", 9).
		Where("c.name", "IS NOT NULL").
		Where("c.age", "BETWEEN ? AND ?", 10, 20).
		Build()

	expectedQuery := `SELECT "c.id", "c.name", "c.telephone" AS "phone", "c.age"
 FROM "customers" AS "c"
 CROSS JOIN "orders" AS "o" ON "o"."customer_id" = "c"."id"
 WHERE ("c.id" = $1) AND ("c.name" IS NOT NULL) AND ("c.age" BETWEEN $2 AND $3)`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9, 10, 20}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestSelectWithInClauseUsingSlice(t *testing.T) {
	c := customer{}

	input := []int{4, 5, 6}
	query, args, _ := Postgres.Select().
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Where("name", "IS NOT NULL").
		Where("id", "in (?,?,?)", input).
		Where("age", "BETWEEN ? AND ?", 10, 20).
		Build()

	expectedQuery := `SELECT "id", "name"
 FROM "customers"
 WHERE ("name" IS NOT NULL) AND ("id" in ($1,$2,$3)) AND ("age" BETWEEN $4 AND $5)`
	if query != expectedQuery {
		t.Errorf("bad query: |%s|", query)
	}

	expectedArgs := []interface{}{4, 5, 6, 10, 20}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}
