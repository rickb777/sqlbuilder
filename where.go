package sqlbuilder

import (
	"reflect"
	"strings"
)

type where struct {
	col, sql string
	args     []interface{}
}

func quoteStrings(s []string, dialect Dialect) []string {
	var quoted []string
	for _, o := range s {
		quoted = append(quoted, dialect.Quote(o))
	}
	return quoted
}

func buildWhereClause(query string, args []interface{}, idx int, wheres []where, dialect Dialect) (string, []interface{}, int) {
	if len(wheres) > 0 {
		var sqls []string

		for _, where := range wheres {
			sql := "(" + dialect.Quote(where.col) + " " + where.sql + ")"
			for _, arg := range where.args {
				value := reflect.ValueOf(arg)
				switch value.Kind() {
				case reflect.Array, reflect.Slice:
					for j := 0; j < value.Len(); j++ {
						p := dialect.Placeholder(idx)
						idx++
						sql = strings.Replace(sql, "?", p, 1)
						args = append(args, value.Index(j).Interface())
					}

				default:
					p := dialect.Placeholder(idx)
					idx++
					sql = strings.Replace(sql, "?", p, 1)
					args = append(args, arg)
				}
			}
			sqls = append(sqls, sql)
		}

		query += "\n WHERE " + strings.Join(sqls, " AND ")
	}
	return query, args, idx
}
