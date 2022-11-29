package predicate

import "github.com/uptrace/bun"

type Product func(query *bun.SelectQuery) *bun.SelectQuery
