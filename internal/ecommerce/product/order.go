package product

import (
	"fmt"
	"github.com/uptrace/bun"
)

type OrderFunc func(query *bun.SelectQuery) *bun.SelectQuery

func Asc(fields ...string) OrderFunc {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, f := range fields {
			q.OrderExpr(fmt.Sprintf("%s ASC", f))
		}
		return q
	}
}

func Desc(fields ...string) OrderFunc {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, f := range fields {
			q.OrderExpr(fmt.Sprintf("%s DESC", f))
		}
		return q
	}
}
