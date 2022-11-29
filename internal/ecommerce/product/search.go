package product

import (
	"context"
	"github.com/huyhvq/ecommerce/internal/ecommerce/model"
	"github.com/huyhvq/ecommerce/internal/ecommerce/predicate"
	"github.com/uptrace/bun"
)

type SearchQuery struct {
	db *bun.DB

	limit  int
	offset int

	order      []OrderFunc
	predicates []predicate.Product
}

func NewSearchQuery(db *bun.DB) *SearchQuery {
	return &SearchQuery{
		db:     db,
		limit:  10,
		offset: 0,
	}
}

type Facet struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Count uint32 `json:"count"`
}

type FacetMap map[string][]*Facet

func (sq *SearchQuery) GetFacets(ctx context.Context) (FacetMap, error) {
	q := sq.db.NewSelect().Column("tsv").Model((*model.Product)(nil))
	sq.filter(q)
	var facets []*Facet
	if err := sq.db.NewSelect().
		ColumnExpr("split_part(word, ':', 1) AS key").
		ColumnExpr("split_part(word, ':', 2) AS value").
		ColumnExpr("ndoc AS count").
		ColumnExpr("row_number() OVER (PARTITION BY split_part(word, ':', 1) ORDER BY ndoc DESC) AS _rank").
		TableExpr("ts_stat($$ ? $$)", q).
		OrderExpr("_rank DESC").
		Scan(ctx, &facets); err != nil {
		return nil, err
	}
	m := make(FacetMap, len(facets))

	for _, facet := range facets {
		m[facet.Key] = append(m[facet.Key], facet)
	}
	return m, nil
}

func (sq *SearchQuery) Search(ctx context.Context, term string) ([]*model.Product, error) {
	var products []*model.Product
	q := sq.db.NewSelect().Model((*model.Product)(nil))
	sq.filter(q)
	q.Where("product_text_search(name, description) @@ plainto_tsquery(?)", term)
	for _, of := range sq.order {
		of(q)
	}
	q.OrderExpr("ts_rank(product_text_search(name, description), plainto_tsquery(?)) DESC", term)
	q.Offset(sq.offset)
	q.Limit(sq.limit)

	if err := q.Scan(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (sq *SearchQuery) filter(q *bun.SelectQuery) {
	for _, p := range sq.predicates {
		p(q)
	}
}

func (sq *SearchQuery) Where(ps ...predicate.Product) *SearchQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

func (sq *SearchQuery) Order(o ...OrderFunc) *SearchQuery {
	sq.order = append(sq.order, o...)
	return sq
}

func (sq *SearchQuery) Limit(limit int) *SearchQuery {
	sq.limit = limit
	return sq
}

func (sq *SearchQuery) Offset(offset int) *SearchQuery {
	sq.offset = offset
	return sq
}
