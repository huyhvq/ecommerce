package product

import (
	"fmt"
	"github.com/huyhvq/ecommerce/internal/ecommerce/predicate"
	"github.com/uptrace/bun"
	"strings"
)

func PriceEQ(price float64) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("price = ?", price)
	}
}

func PriceGT(price float64) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("price > ?", price)
	}
}

func PriceGTE(price float64) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("price >= ?", price)
	}
}

func PriceLT(price float64) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("price < ?", price)
	}
}

func PriceLTE(price float64) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("price <= ?", price)
	}
}

func AttributeIn(attributes ...string) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, attr := range attributes {
			attr = strings.ReplaceAll(attr, ":", "\\:")
			q = q.Where("tsv @@ ?::tsquery", attr)
		}
		return q
	}

}

func AttributeNotIn(attributes ...string) predicate.Product {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, attr := range attributes {
			attr = strings.ReplaceAll(attr, ":", "\\:")
			q = q.Where("tsv @@ ?::tsquery", fmt.Sprintf("! %s", attr))
		}
		return q
	}

}
