package product

import (
	"context"
	"github.com/huyhvq/ecommerce/internal/ecommerce/model"
	"github.com/uptrace/bun"
	"time"
)

type LogStore struct {
	db *bun.DB
}

func NewLogStore(db *bun.DB) *LogStore {
	return &LogStore{
		db: db,
	}
}

type Facet struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Count uint32 `json:"count"`
}

type FacetMap map[string][]*Facet

func (s *LogStore) Create(ctx context.Context, term string, searchCtx []string) error {
	l := &model.Log{
		SearchTerm: term,
		Context:    searchCtx,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if _, err := s.db.NewInsert().Model(l).Exec(ctx); err != nil {
		return err
	}
	return nil
}
