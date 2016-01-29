package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestDeleteWithWhereMySQL(t *testing.T) {
	query, args := MySQLQuoted.Delete().
		From("customers").
		Where("id", "= ?", 9).
		Build()

	expectedQuery := "DELETE FROM `customers`\n WHERE (`id` = ?)"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestDeleteWithWherePostgres(t *testing.T) {
	query, args := Postgres.Delete().
		From("customers").
		Where("id", "= ?", 9).
		Build()

	expectedQuery := `DELETE FROM "customers"
 WHERE ("id" = $1)`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}
