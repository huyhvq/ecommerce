package model

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

type Product struct {
	ID          uint64    `json:"id" bun:",pk,autoincrement"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Attributes  []string  `json:"attributes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*Product) BeforeCreateTable(ctx context.Context, q *bun.CreateTableQuery) error {
	q.ColumnExpr("tsv tsvector")
	return nil
}
