package api

import (
	"context"
	"fmt"
	pls "github.com/huyhvq/ecommerce/internal/ecommerce/log"
	"github.com/huyhvq/ecommerce/internal/ecommerce/predicate"
	"github.com/huyhvq/ecommerce/internal/ecommerce/product"
	"github.com/huyhvq/ecommerce/pkg/request"
	"github.com/huyhvq/ecommerce/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

type SearchParams struct {
	Price      float64 `json:"price"`
	Categories string  `json:"categories"`
	Moods      string  `json:"moods"`
	Pace       string  `json:"pace"`
}

var filterParams = map[string]string{
	"price":      "float",
	"categories": "string",
	"colors":     "string",
	"size":       "string",
	"brand":      "string",
	"collection": "string",
}

var sortParams = map[string]string{
	"sort": "string",
}

func (app *application) productSearch(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("term") == "" {
		app.errorMessage(w, r, http.StatusUnprocessableEntity, "term is required", nil)
		return
	}
	term := r.URL.Query().Get("term")
	sq := product.NewSearchQuery(app.db.DB)
	for _, f := range request.QueryStringParser(r.URL.RawQuery, filterParams) {
		if f.Field == "price" {
			sq.Where(app.filterPrice(f))
			continue
		}
		attributes := strings.Split(f.Value.(string), ",")
		for i, attr := range attributes {
			attributes[i] = fmt.Sprintf("%s:%s", f.Field, attr)
		}
		sq.Where(product.AttributeIn(attributes...))
	}

	for _, f := range request.QueryStringParser(r.URL.RawQuery, sortParams) {
		if f.Operator == "asc" {
			sq.Order(product.Asc(f.Value.(string)))
		}
		if f.Operator == "desc" {
			sq.Order(product.Desc(f.Value.(string)))
		}
	}

	if p := r.URL.Query().Get("limit"); p != "" {
		limit, err := strconv.Atoi(p)
		if err != nil {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "limit is invalid", nil)
			return
		}
		sq.Limit(limit)
	}

	if p := r.URL.Query().Get("offset"); p != "" {
		offset, err := strconv.Atoi(p)
		if err != nil {
			app.errorMessage(w, r, http.StatusUnprocessableEntity, "offset is invalid", nil)
			return
		}
		sq.Offset(offset)
	}

	ls := pls.NewLogStore(app.db.DB)
	searchCond := make([]string, 0)
	for k, vs := range r.URL.Query() {
		if len(vs) == 1 {
			searchCond = append(searchCond, fmt.Sprintf("%s:%s", k, vs[0]))
		} else {
			for _, v := range vs {
				searchCond = append(searchCond, fmt.Sprintf("%s:%s", k, v))
			}
		}
	}
	go ls.Create(context.Background(), term, searchCond)

	f, err := sq.Order(product.Desc("price")).Search(r.Context(), term)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, f); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) productFacets(w http.ResponseWriter, r *http.Request) {
	sq := product.NewSearchQuery(app.db.DB)
	f, err := sq.GetFacets(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, f); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) filterPrice(filter request.FilteredResult) predicate.Product {
	v := filter.Value.(float64)
	switch strings.Trim(filter.Operator, " ") {
	case "<=":
		return product.PriceLTE(v)
	case ">=":
		return product.PriceGTE(v)
	case "<":
		return product.PriceLT(v)
	case ">":
		return product.PriceGT(v)
	default:
		return product.PriceEQ(v)
	}
}
